package controller

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/thanhpk/randstr"
)

const PaymentMethodWechat = "wechat"

// 微信支付统一下单请求
type WechatUnifiedOrderReq struct {
	XMLName        xml.Name `xml:"xml"`
	AppId          string   `xml:"appid"`
	MchId          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	Body           string   `xml:"body"`
	OutTradeNo     string   `xml:"out_trade_no"`
	TotalFee       int      `xml:"total_fee"`
	SpbillCreateIp string   `xml:"spbill_create_ip"`
	NotifyUrl      string   `xml:"notify_url"`
	TradeType      string   `xml:"trade_type"`
}

// 微信支付统一下单响应
type WechatUnifiedOrderResp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	CodeUrl    string `xml:"code_url"`
	PrepayId   string `xml:"prepay_id"`
}

// 微信支付回调通知
type WechatNotify struct {
	ReturnCode    string `xml:"return_code"`
	ReturnMsg     string `xml:"return_msg"`
	ResultCode    string `xml:"result_code"`
	OutTradeNo    string `xml:"out_trade_no"`
	TransactionId string `xml:"transaction_id"`
	TimeEnd       string `xml:"time_end"`
	TotalFee      int    `xml:"total_fee"`
	Sign          string `xml:"sign"`
}

// 请求微信支付
func RequestWechatPay(c *gin.Context) {
	var req struct {
		Amount        int64  `json:"amount"`
		PaymentMethod string `json:"payment_method"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "参数错误"})
		return
	}

	if req.PaymentMethod != PaymentMethodWechat {
		c.JSON(200, gin.H{"message": "error", "data": "不支持的支付渠道"})
		return
	}

	minTopup := getWechatMinTopup()
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

	payMoney := getWechatPayMoney(float64(req.Amount), group)
	if payMoney < 0.01 {
		c.JSON(200, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}

	// 转换为分
	totalFee := int(payMoney * 100)

	// 生成订单号
	reference := fmt.Sprintf("wechat-ref-%d-%d-%s", user.Id, time.Now().UnixMilli(), randstr.String(4))
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
		PaymentMethod: PaymentMethodWechat,
		CreateTime:    time.Now().Unix(),
		Status:        common.TopUpStatusPending,
	}

	if err := topUp.Insert(); err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "创建订单失败"})
		return
	}

	// 调用微信统一下单
	notifyUrl := system_setting.ServerAddress + "/api/wechat/notify"
	codeUrl, err := wechatUnifiedOrder(referenceId, totalFee, notifyUrl, c.ClientIP())
	if err != nil {
		log.Printf("微信统一下单失败: %v", err)
		c.JSON(200, gin.H{"message": "error", "data": "拉起支付失败"})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data": gin.H{
			"pay_url":  codeUrl,
			"order_id": referenceId,
		},
	})
}

// 微信统一下单
func wechatUnifiedOrder(outTradeNo string, totalFee int, notifyUrl string, clientIP string) (string, error) {
	if setting.WechatMchId == "" || setting.WechatApiV2Key == "" {
		return "", fmt.Errorf("微信支付未配置")
	}

	// 构建请求参数
	params := map[string]string{
		"mch_id":           setting.WechatMchId,
		"nonce_str":        randstr.String(32),
		"body":             "AI服务充值",
		"out_trade_no":     outTradeNo,
		"total_fee":        fmt.Sprintf("%d", totalFee),
		"spbill_create_ip": clientIP,
		"notify_url":       notifyUrl,
		"trade_type":       "NATIVE",
	}

	// 如果配置了AppID则添加
	if setting.WechatAppId != "" {
		params["appid"] = setting.WechatAppId
	}

	// 生成签名
	sign := wechatSign(params, setting.WechatApiV2Key)
	params["sign"] = sign

	// 构建XML请求
	reqXML := buildWechatXML(params)

	// 发送请求
	resp, err := http.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", "application/xml", strings.NewReader(reqXML))
	if err != nil {
		return "", fmt.Errorf("请求微信API失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result WechatUnifiedOrderResp
	if err := xml.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.ReturnCode != "SUCCESS" {
		return "", fmt.Errorf("微信返回错误: %s", result.ReturnMsg)
	}

	if result.ResultCode != "SUCCESS" {
		return "", fmt.Errorf("下单失败: %s - %s", result.ErrCode, result.ErrCodeDes)
	}

	return result.CodeUrl, nil
}

// 微信支付签名（MD5）
func wechatSign(params map[string]string, apiKey string) string {
	// 1. 参数排序
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" && params[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 2. 拼接字符串
	var signStr strings.Builder
	for i, k := range keys {
		if i > 0 {
			signStr.WriteString("&")
		}
		signStr.WriteString(k)
		signStr.WriteString("=")
		signStr.WriteString(params[k])
	}
	signStr.WriteString("&key=")
	signStr.WriteString(apiKey)

	// 3. MD5加密并转大写
	hash := md5.Sum([]byte(signStr.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// 验证微信签名
func verifyWechatSign(params map[string]string, apiKey string) bool {
	receivedSign := params["sign"]
	delete(params, "sign")

	calculatedSign := wechatSign(params, apiKey)
	return receivedSign == calculatedSign
}

// 微信支付回调
func WechatNotifyHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("微信回调读取数据失败: %v", err)
		wechatNotifyResponse(c, "FAIL", "读取数据失败")
		return
	}

	var notify WechatNotify
	if err := xml.Unmarshal(body, &notify); err != nil {
		log.Printf("微信回调解析数据失败: %v", err)
		wechatNotifyResponse(c, "FAIL", "解析数据失败")
		return
	}

	// 将XML解析为map用于签名验证
	params := make(map[string]string)
	decoder := xml.NewDecoder(strings.NewReader(string(body)))
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("微信回调XML解析错误: %v", err)
			wechatNotifyResponse(c, "FAIL", "签名验证失败")
			return
		}

		if se, ok := token.(xml.StartElement); ok {
			var value string
			if err := decoder.DecodeElement(&value, &se); err == nil {
				params[se.Name.Local] = value
			}
		}
	}

	// 验签
	if !verifyWechatSign(params, setting.WechatApiV2Key) {
		log.Printf("微信回调签名验证失败")
		wechatNotifyResponse(c, "FAIL", "签名验证失败")
		return
	}

	if notify.ReturnCode != "SUCCESS" {
		log.Printf("微信回调返回失败: %s", notify.ReturnMsg)
		wechatNotifyResponse(c, "FAIL", notify.ReturnMsg)
		return
	}

	if notify.ResultCode != "SUCCESS" {
		log.Printf("微信支付失败: %s", notify.OutTradeNo)
		wechatNotifyResponse(c, "FAIL", "支付失败")
		return
	}

	// 查询订单
	LockOrder(notify.OutTradeNo)
	defer UnlockOrder(notify.OutTradeNo)

	topUp := model.GetTopUpByTradeNo(notify.OutTradeNo)
	if topUp == nil {
		log.Printf("微信回调: 订单不存在 - %s", notify.OutTradeNo)
		wechatNotifyResponse(c, "FAIL", "订单不存在")
		return
	}

	// 幂等性检查
	if topUp.Status != common.TopUpStatusPending {
		log.Printf("微信回调: 订单已处理 - %s, 状态: %s", notify.OutTradeNo, topUp.Status)
		wechatNotifyResponse(c, "SUCCESS", "OK")
		return
	}

	// 处理充值
	dAmount := decimal.NewFromInt(topUp.Amount)
	dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
	quotaToAdd := int(dAmount.Mul(dQuotaPerUnit).IntPart())

	err = model.IncreaseUserQuota(topUp.UserId, quotaToAdd, true)
	if err != nil {
		log.Printf("微信充值增加配额失败: %v", err)
		wechatNotifyResponse(c, "FAIL", "处理失败")
		return
	}

	// 更新订单状态
	topUp.Status = common.TopUpStatusSuccess
	topUp.CompleteTime = time.Now().Unix()
	if err := topUp.Update(); err != nil {
		log.Printf("微信充值更新订单失败: %v", err)
		wechatNotifyResponse(c, "FAIL", "更新订单失败")
		return
	}

	// 记录日志
	model.RecordLog(topUp.UserId, model.LogTypeTopup, fmt.Sprintf("使用微信支付成功，充值金额: %d，支付金额：%.2f", quotaToAdd, topUp.Money))

	log.Printf("微信支付成功 - 订单号: %s, 微信流水号: %s, 金额: %.2f, 配额: %d",
		notify.OutTradeNo, notify.TransactionId, topUp.Money, quotaToAdd)

	wechatNotifyResponse(c, "SUCCESS", "OK")
}

// 微信回调响应
func wechatNotifyResponse(c *gin.Context, code, msg string) {
	resp := fmt.Sprintf("<xml><return_code><![CDATA[%s]]></return_code><return_msg><![CDATA[%s]]></return_msg></xml>", code, msg)
	c.String(200, resp)
}

// 构建微信XML请求
func buildWechatXML(params map[string]string) string {
	var xml strings.Builder
	xml.WriteString("<xml>")
	for k, v := range params {
		xml.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	xml.WriteString("</xml>")
	return xml.String()
}

// 获取微信支付金额
func getWechatPayMoney(amount float64, group string) float64 {
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

// 获取微信最小充值
func getWechatMinTopup() int64 {
	minTopup := setting.WechatMinTopUp
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		dMinTopup := decimal.NewFromInt(int64(minTopup))
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		minTopup = int(dMinTopup.Mul(dQuotaPerUnit).IntPart())
	}
	return int64(minTopup)
}

// 获取微信支付金额（用于前端显示）
func RequestWechatAmount(c *gin.Context) {
	var req struct {
		Amount int64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "参数错误"})
		return
	}

	minTopup := getWechatMinTopup()
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

	payMoney := getWechatPayMoney(float64(req.Amount), group)
	if payMoney <= 0.01 {
		c.JSON(200, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}

	c.JSON(200, gin.H{"message": "success", "data": fmt.Sprintf("%.2f", payMoney)})
}
