package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/smartwalle/alipay/v3"
	"github.com/thanhpk/randstr"
)

const PaymentMethodAlipay = "alipay"

var alipayClient *alipay.Client

// 初始化支付宝客户端
func InitAlipayClient() error {
	if setting.AlipayAppId == "" || setting.AlipayPrivateKey == "" {
		return fmt.Errorf("支付宝配置不完整")
	}

	var err error
	alipayClient, err = alipay.New(
		setting.AlipayAppId,
		setting.AlipayPrivateKey,
		false, // 是否使用沙箱环境
	)
	if err != nil {
		return err
	}

	// 加载支付宝公钥
	if setting.AlipayPublicKey != "" {
		if err := alipayClient.LoadAliPayPublicKey(setting.AlipayPublicKey); err != nil {
			return err
		}
	}

	log.Printf("支付宝客户端初始化成功")
	return nil
}

// 请求支付宝支付
func RequestAlipayPay(c *gin.Context) {
	var req struct {
		Amount        int64  `json:"amount"`
		PaymentMethod string `json:"payment_method"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "参数错误"})
		return
	}

	if req.PaymentMethod != PaymentMethodAlipay {
		c.JSON(200, gin.H{"message": "error", "data": "不支持的支付渠道"})
		return
	}

	minTopup := getAlipayMinTopup()
	if req.Amount < minTopup {
		c.JSON(200, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", minTopup)})
		return
	}

	id := c.GetInt("id")
	user, _ := model.GetUserById(id, false)

	// 计算实际支付金额
	group, err := model.GetUserGroup(id, true)
	if err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "获取用户分组失败"})
		return
	}

	payMoney := getAlipayPayMoney(float64(req.Amount), group)
	if payMoney < 0.01 {
		c.JSON(200, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}

	// 生成订单号
	reference := fmt.Sprintf("alipay-ref-%d-%d-%s", user.Id, time.Now().UnixMilli(), randstr.String(4))
	referenceId := "ref_" + common.Sha1([]byte(reference))

	// 创建订单记录
	amount := req.Amount
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		dAmount := decimal.NewFromInt(amount)
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		amount = dAmount.Div(dQuotaPerUnit).IntPart()
	}

	topUp := &model.TopUp{
		UserId:        id,
		Amount:        amount,
		Money:         payMoney,
		TradeNo:       referenceId,
		PaymentMethod: PaymentMethodAlipay,
		CreateTime:    time.Now().Unix(),
		Status:        common.TopUpStatusPending,
	}

	if err := topUp.Insert(); err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "创建订单失败"})
		return
	}

	// 调用支付宝预下单（扫码支付）
	notifyUrl := system_setting.ServerAddress + "/api/alipay/notify"
	qrCode, err := alipayQRCodePay(referenceId, payMoney, notifyUrl)
	if err != nil {
		log.Printf("支付宝预下单失败: %v", err)
		c.JSON(200, gin.H{"message": "error", "data": "拉起支付失败"})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data": gin.H{
			"pay_url":  qrCode,
			"order_id": referenceId,
		},
	})
}

// 支付宝扫码支付（当面付）
func alipayQRCodePay(outTradeNo string, totalAmount float64, notifyUrl string) (string, error) {
	// 懒加载初始化客户端
	if alipayClient == nil {
		if err := InitAlipayClient(); err != nil {
			return "", err
		}
	}

	var p = alipay.TradePreCreate{}
	p.OutTradeNo = outTradeNo
	p.Subject = "AI服务充值"
	p.TotalAmount = fmt.Sprintf("%.2f", totalAmount)
	p.NotifyURL = notifyUrl

	rsp, err := alipayClient.TradePreCreate(context.Background(), p)
	if err != nil {
		return "", fmt.Errorf("调用支付宝API失败: %w", err)
	}

	if rsp.IsFailure() {
		return "", fmt.Errorf("支付宝返回错误: %s - %s", rsp.Code, rsp.SubMsg)
	}

	if rsp.QRCode == "" {
		return "", fmt.Errorf("支付宝未返回二维码")
	}

	return rsp.QRCode, nil
}

// 支付宝支付回调
func AlipayNotifyHandler(c *gin.Context) {
	// 解析表单数据
	if err := c.Request.ParseForm(); err != nil {
		log.Printf("支付宝回调解析表单失败: %v", err)
		c.String(200, "fail")
		return
	}

	// 懒加载初始化客户端
	if alipayClient == nil {
		if err := InitAlipayClient(); err != nil {
			log.Printf("支付宝回调初始化客户端失败: %v", err)
			c.String(200, "fail")
			return
		}
	}

	// 验签
	notification, err := alipayClient.DecodeNotification(c.Request.Form)
	if err != nil {
		log.Printf("支付宝签名验证失败: %v", err)
		c.String(200, "fail")
		return
	}

	// 获取通知参数
	outTradeNo := notification.OutTradeNo
	tradeStatus := notification.TradeStatus
	tradeNo := notification.TradeNo

	log.Printf("支付宝回调 - 订单号: %s, 支付宝流水号: %s, 状态: %s", outTradeNo, tradeNo, tradeStatus)

	// 查询订单
	LockOrder(outTradeNo)
	defer UnlockOrder(outTradeNo)

	topUp := model.GetTopUpByTradeNo(outTradeNo)
	if topUp == nil {
		log.Printf("支付宝回调: 订单不存在 - %s", outTradeNo)
		c.String(200, "fail")
		return
	}

	// 幂等性检查
	if topUp.Status != common.TopUpStatusPending {
		log.Printf("支付宝回调: 订单已处理 - %s, 状态: %s", outTradeNo, topUp.Status)
		c.String(200, "success")
		return
	}

	// 交易成功
	if tradeStatus == "TRADE_SUCCESS" || tradeStatus == "TRADE_FINISHED" {
		// 增加用户配额
		dAmount := decimal.NewFromInt(topUp.Amount)
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		quotaToAdd := int(dAmount.Mul(dQuotaPerUnit).IntPart())

		err := model.IncreaseUserQuota(topUp.UserId, quotaToAdd, true)
		if err != nil {
			log.Printf("支付宝充值增加配额失败: %v", err)
			c.String(200, "fail")
			return
		}

		// 更新订单状态
		topUp.Status = common.TopUpStatusSuccess
		topUp.CompleteTime = time.Now().Unix()
		if err := topUp.Update(); err != nil {
			log.Printf("支付宝充值更新订单失败: %v", err)
			c.String(200, "fail")
			return
		}

		// 记录日志
		model.RecordLog(topUp.UserId, model.LogTypeTopup, fmt.Sprintf("使用支付宝支付成功，充值金额: %d，支付金额：%.2f", quotaToAdd, topUp.Money))

		log.Printf("支付宝支付成功 - 订单号: %s, 支付宝流水号: %s, 金额: %.2f, 配额: %d",
			outTradeNo, tradeNo, topUp.Money, quotaToAdd)

		c.String(200, "success")
	} else {
		log.Printf("支付宝回调状态异常: %s - %s", outTradeNo, tradeStatus)
		c.String(200, "fail")
	}
}

// 获取支付宝支付金额
func getAlipayPayMoney(amount float64, group string) float64 {
	originalAmount := amount
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		dAmount := decimal.NewFromFloat(amount)
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		amount = dAmount.Div(dQuotaPerUnit).InexactFloat64()
	}

	topupGroupRatio := common.GetTopupGroupRatio(group)
	if topupGroupRatio == 0 {
		topupGroupRatio = 1
	}

	// 应用折扣
	discount := 1.0
	if ds, ok := operation_setting.GetPaymentSetting().AmountDiscount[int(originalAmount)]; ok {
		if ds > 0 {
			discount = ds
		}
	}

	dAmount := decimal.NewFromFloat(amount)
	dPrice := decimal.NewFromFloat(operation_setting.Price)
	dGroupRatio := decimal.NewFromFloat(topupGroupRatio)
	dDiscount := decimal.NewFromFloat(discount)

	payMoney := dAmount.Mul(dPrice).Mul(dGroupRatio).Mul(dDiscount)

	return payMoney.InexactFloat64()
}

// 获取支付宝最小充值
func getAlipayMinTopup() int64 {
	minTopup := setting.AlipayMinTopUp
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		dMinTopup := decimal.NewFromInt(int64(minTopup))
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		minTopup = int(dMinTopup.Mul(dQuotaPerUnit).IntPart())
	}
	return int64(minTopup)
}

// 获取支付宝支付金额（用于前端显示）
func RequestAlipayAmount(c *gin.Context) {
	var req struct {
		Amount int64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "参数错误"})
		return
	}

	minTopup := getAlipayMinTopup()
	if req.Amount < minTopup {
		c.JSON(200, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", minTopup)})
		return
	}

	id := c.GetInt("id")
	group, err := model.GetUserGroup(id, true)
	if err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "获取用户分组失败"})
		return
	}

	payMoney := getAlipayPayMoney(float64(req.Amount), group)
	if payMoney <= 0.01 {
		c.JSON(200, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}

	c.JSON(200, gin.H{"message": "success", "data": fmt.Sprintf("%.2f", payMoney)})
}
