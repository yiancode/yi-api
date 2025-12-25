# 微信支付和支付宝支付实现方案

## 一、总体架构设计

### 1.1 技术选型

| 支付方式 | SDK/API | 支付场景 |
|---------|---------|---------|
| **微信支付** | 微信支付 API v2 | Native扫码支付 |
| **支付宝** | 支付宝开放平台 SDK | 当面付（扫码支付） |

### 1.2 配置参数

**微信支付配置**：
- `WechatMchId` - 商户号
- `WechatApiV2Key` - API v2密钥

**支付宝配置**：
- `AlipayAppId` - 应用APPID
- `AlipayPrivateKey` - 应用私钥
- `AlipayPublicKey` - 支付宝公钥

### 1.3 文件结构

```
后端：
├── controller/
│   ├── topup_wechat.go         # 微信支付控制器
│   └── topup_alipay.go         # 支付宝支付控制器
├── setting/
│   ├── payment_wechat.go       # 微信支付配置
│   └── payment_alipay.go       # 支付宝配置
└── router/
    └── api-router.go           # 路由配置（新增回调路由）

前端：
├── web/src/pages/Setting/Payment/
│   ├── SettingsPaymentGatewayWechat.jsx    # 微信支付配置页
│   └── SettingsPaymentGatewayAlipay.jsx    # 支付宝配置页
└── web/src/components/topup/
    └── index.jsx               # 充值页面（添加微信/支付宝支付按钮）
```

---

## 二、后端实现

### 2.1 配置模块

#### `setting/payment_wechat.go`

```go
package setting

var WechatMchId = ""          // 微信商户号
var WechatApiV2Key = ""       // 微信API v2密钥
var WechatMinTopUp = 1        // 最小充值金额
```

#### `setting/payment_alipay.go`

```go
package setting

var AlipayAppId = ""          // 支付宝AppID
var AlipayPrivateKey = ""     // 应用私钥
var AlipayPublicKey = ""      // 支付宝公钥
var AlipayMinTopUp = 1        // 最小充值金额
```

---

### 2.2 微信支付实现

#### `controller/topup_wechat.go` - 核心代码结构

```go
package controller

import (
    "crypto/md5"
    "encoding/hex"
    "encoding/xml"
    "fmt"
    "io"
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

    if req.Amount < int64(setting.WechatMinTopUp) {
        c.JSON(200, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", setting.WechatMinTopUp)})
        return
    }

    id := c.GetInt("id")
    user, _ := model.GetUserById(id, false)

    // 计算实际支付金额（分）
    group, err := model.GetUserGroup(id, true)
    if err != nil {
        c.JSON(200, gin.H{"message": "error", "data": "获取用户分组失败"})
        return
    }

    payMoney := getWechatPayMoney(float64(req.Amount), group)
    totalFee := int(payMoney * 100) // 转换为分

    // 生成订单号
    reference := fmt.Sprintf("wechat-ref-%d-%d-%s", user.Id, time.Now().UnixMilli(), randstr.String(4))
    referenceId := "ref_" + common.Sha1([]byte(reference))

    // 创建订单
    topUp := &model.TopUp{
        UserId:        id,
        Amount:        req.Amount,
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
            "pay_url": codeUrl,
            "order_id": referenceId,
        },
    })
}

// 微信统一下单
func wechatUnifiedOrder(outTradeNo string, totalFee int, notifyUrl string, clientIP string) (string, error) {
    // 构建请求参数
    params := map[string]string{
        "appid":            setting.WechatAppId,
        "mch_id":           setting.WechatMchId,
        "nonce_str":        randstr.String(32),
        "body":             "AI服务充值",
        "out_trade_no":     outTradeNo,
        "total_fee":        fmt.Sprintf("%d", totalFee),
        "spbill_create_ip": clientIP,
        "notify_url":       notifyUrl,
        "trade_type":       "NATIVE",
    }

    // 生成签名
    sign := wechatSign(params, setting.WechatApiV2Key)
    params["sign"] = sign

    // 构建XML请求
    reqXML := buildWechatXML(params)

    // 发送请求
    resp, err := http.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", "application/xml", strings.NewReader(reqXML))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    var result WechatUnifiedOrderResp
    if err := xml.Unmarshal(body, &result); err != nil {
        return "", err
    }

    if result.ReturnCode != "SUCCESS" {
        return "", fmt.Errorf("微信返回错误: %s", result.ReturnMsg)
    }

    if result.ResultCode != "SUCCESS" {
        return "", fmt.Errorf("下单失败: %s", result.ErrCodeDes)
    }

    return result.CodeUrl, nil
}

// 微信支付签名
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
    for _, k := range keys {
        signStr.WriteString(k)
        signStr.WriteString("=")
        signStr.WriteString(params[k])
        signStr.WriteString("&")
    }
    signStr.WriteString("key=")
    signStr.WriteString(apiKey)

    // 3. MD5加密并转大写
    hash := md5.Sum([]byte(signStr.String()))
    return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// 微信支付回调
func WechatNotifyHandler(c *gin.Context) {
    body, err := io.ReadAll(c.Request.Body)
    if err != nil {
        wechatNotifyResponse(c, "FAIL", "读取数据失败")
        return
    }

    var notify WechatNotify
    if err := xml.Unmarshal(body, &notify); err != nil {
        wechatNotifyResponse(c, "FAIL", "解析数据失败")
        return
    }

    // 验签
    if !verifyWechatSign(body, setting.WechatApiV2Key) {
        wechatNotifyResponse(c, "FAIL", "签名验证失败")
        return
    }

    // 查询订单
    LockOrder(notify.OutTradeNo)
    defer UnlockOrder(notify.OutTradeNo)

    topUp := model.GetTopUpByTradeNo(notify.OutTradeNo)
    if topUp == nil {
        wechatNotifyResponse(c, "FAIL", "订单不存在")
        return
    }

    if topUp.Status != common.TopUpStatusPending {
        wechatNotifyResponse(c, "SUCCESS", "OK")
        return
    }

    // 处理充值
    if notify.ResultCode == "SUCCESS" {
        err := model.Recharge(notify.OutTradeNo, "")
        if err != nil {
            log.Printf("微信充值处理失败: %v", err)
            wechatNotifyResponse(c, "FAIL", "处理失败")
            return
        }

        log.Printf("微信支付成功 - 订单号: %s, 金额: %.2f", notify.OutTradeNo, topUp.Money)
        wechatNotifyResponse(c, "SUCCESS", "OK")
    } else {
        wechatNotifyResponse(c, "FAIL", "支付失败")
    }
}

// 微信回调响应
func wechatNotifyResponse(c *gin.Context, code, msg string) {
    resp := fmt.Sprintf("<xml><return_code><![CDATA[%s]]></return_code><return_msg><![CDATA[%s]]></return_msg></xml>", code, msg)
    c.String(200, resp)
}

// 验证微信签名
func verifyWechatSign(xmlData []byte, apiKey string) bool {
    // 解析XML获取参数
    params := make(map[string]string)
    // ... XML解析逻辑

    receivedSign := params["sign"]
    delete(params, "sign")

    calculatedSign := wechatSign(params, apiKey)
    return receivedSign == calculatedSign
}

func getWechatPayMoney(amount float64, group string) float64 {
    if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
        amount = amount / common.QuotaPerUnit
    }

    topupGroupRatio := common.GetTopupGroupRatio(group)
    if topupGroupRatio == 0 {
        topupGroupRatio = 1
    }

    discount := 1.0
    if ds, ok := operation_setting.GetPaymentSetting().AmountDiscount[int(amount)]; ok {
        if ds > 0 {
            discount = ds
        }
    }

    return amount * operation_setting.Price * topupGroupRatio * discount
}

func buildWechatXML(params map[string]string) string {
    var xml strings.Builder
    xml.WriteString("<xml>")
    for k, v := range params {
        xml.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
    }
    xml.WriteString("</xml>")
    return xml.String()
}
```

---

### 2.3 支付宝支付实现

#### `controller/topup_alipay.go` - 核心代码结构

```go
package controller

import (
    "encoding/json"
    "fmt"
    "net/url"
    "time"

    "github.com/QuantumNous/new-api/common"
    "github.com/QuantumNous/new-api/model"
    "github.com/QuantumNous/new-api/setting"
    "github.com/QuantumNous/new-api/setting/operation_setting"
    "github.com/QuantumNous/new-api/setting/system_setting"
    "github.com/gin-gonic/gin"
    "github.com/smartwalle/alipay/v3"
    "github.com/thanhpk/randstr"
)

const PaymentMethodAlipay = "alipay"

var alipayClient *alipay.Client

// 初始化支付宝客户端
func InitAlipayClient() error {
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
    if err := alipayClient.LoadAliPayPublicKey(setting.AlipayPublicKey); err != nil {
        return err
    }

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

    if req.Amount < int64(setting.AlipayMinTopUp) {
        c.JSON(200, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", setting.AlipayMinTopUp)})
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

    // 生成订单号
    reference := fmt.Sprintf("alipay-ref-%d-%d-%s", user.Id, time.Now().UnixMilli(), randstr.String(4))
    referenceId := "ref_" + common.Sha1([]byte(reference))

    // 创建订单
    topUp := &model.TopUp{
        UserId:        id,
        Amount:        req.Amount,
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
            "pay_url": qrCode,
            "order_id": referenceId,
        },
    })
}

// 支付宝扫码支付（当面付）
func alipayQRCodePay(outTradeNo string, totalAmount float64, notifyUrl string) (string, error) {
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

    rsp, err := alipayClient.TradePreCreate(p)
    if err != nil {
        return "", err
    }

    if rsp.Code != alipay.CodeSuccess {
        return "", fmt.Errorf("支付宝返回错误: %s - %s", rsp.Code, rsp.Msg)
    }

    return rsp.QRCode, nil
}

// 支付宝支付回调
func AlipayNotifyHandler(c *gin.Context) {
    // 解析表单数据
    c.Request.ParseForm()

    // 验签
    if alipayClient == nil {
        InitAlipayClient()
    }

    ok, err := alipayClient.VerifySign(c.Request.Form)
    if err != nil || !ok {
        log.Printf("支付宝签名验证失败: %v", err)
        c.String(200, "fail")
        return
    }

    // 获取通知参数
    outTradeNo := c.Request.FormValue("out_trade_no")
    tradeStatus := c.Request.FormValue("trade_status")
    tradeNo := c.Request.FormValue("trade_no")

    // 查询订单
    LockOrder(outTradeNo)
    defer UnlockOrder(outTradeNo)

    topUp := model.GetTopUpByTradeNo(outTradeNo)
    if topUp == nil {
        log.Printf("支付宝回调: 订单不存在 - %s", outTradeNo)
        c.String(200, "fail")
        return
    }

    if topUp.Status != common.TopUpStatusPending {
        c.String(200, "success")
        return
    }

    // 交易成功
    if tradeStatus == "TRADE_SUCCESS" || tradeStatus == "TRADE_FINISHED" {
        err := model.Recharge(outTradeNo, "")
        if err != nil {
            log.Printf("支付宝充值处理失败: %v", err)
            c.String(200, "fail")
            return
        }

        log.Printf("支付宝支付成功 - 订单号: %s, 支付宝流水号: %s, 金额: %.2f",
            outTradeNo, tradeNo, topUp.Money)
        c.String(200, "success")
    } else {
        c.String(200, "fail")
    }
}

func getAlipayPayMoney(amount float64, group string) float64 {
    if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
        amount = amount / common.QuotaPerUnit
    }

    topupGroupRatio := common.GetTopupGroupRatio(group)
    if topupGroupRatio == 0 {
        topupGroupRatio = 1
    }

    discount := 1.0
    if ds, ok := operation_setting.GetPaymentSetting().AmountDiscount[int(amount)]; ok {
        if ds > 0 {
            discount = ds
        }
    }

    return amount * operation_setting.Price * topupGroupRatio * discount
}
```

---

## 三、路由配置

### 在 `router/api-router.go` 中添加：

```go
// 微信支付
apiRouter.POST("/wechat/webhook", controller.WechatNotifyHandler)
apiRouter.POST("/wechat/notify", controller.WechatNotifyHandler)

// 支付宝
apiRouter.POST("/alipay/webhook", controller.AlipayNotifyHandler)
apiRouter.POST("/alipay/notify", controller.AlipayNotifyHandler)

// 用户支付路由（在 selfRoute 中添加）
selfRoute.POST("/wechat/pay", middleware.CriticalRateLimit(), controller.RequestWechatPay)
selfRoute.POST("/alipay/pay", middleware.CriticalRateLimit(), controller.RequestAlipayPay)
```

---

## 四、前端实现

### 4.1 配置页面

#### `web/src/pages/Setting/Payment/SettingsPaymentGatewayWechat.jsx`

```jsx
import React, { useState, useEffect } from 'react';
import { Form, Input, Button, Toast } from '@douyinfe/semi-ui';
import { API } from '../../../helpers';

export default function SettingsPaymentGatewayWechat() {
  const [loading, setLoading] = useState(false);
  const [config, setConfig] = useState({
    wechat_mch_id: '',
    wechat_api_v2_key: '',
    wechat_min_topup: 1,
  });

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    try {
      const res = await API.get('/api/option');
      if (res.data.success) {
        setConfig({
          wechat_mch_id: res.data.data.WechatMchId || '',
          wechat_api_v2_key: res.data.data.WechatApiV2Key || '',
          wechat_min_topup: res.data.data.WechatMinTopUp || 1,
        });
      }
    } catch (error) {
      Toast.error('加载配置失败');
    }
  };

  const handleSubmit = async (values) => {
    setLoading(true);
    try {
      const res = await API.put('/api/option', {
        WechatMchId: values.wechat_mch_id,
        WechatApiV2Key: values.wechat_api_v2_key,
        WechatMinTopUp: parseInt(values.wechat_min_topup),
      });

      if (res.data.success) {
        Toast.success('保存成功');
      } else {
        Toast.error(res.data.message || '保存失败');
      }
    } catch (error) {
      Toast.error('保存失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-6">
      <h2 className="text-xl font-bold mb-4">微信支付配置</h2>
      <Form
        initValues={config}
        onSubmit={handleSubmit}
        style={{ maxWidth: 600 }}
      >
        <Form.Input
          field="wechat_mch_id"
          label="商户号"
          placeholder="请输入微信支付商户号"
          rules={[{ required: true, message: '请输入商户号' }]}
        />
        <Form.Input
          field="wechat_api_v2_key"
          label="API v2密钥"
          placeholder="请输入微信支付API v2密钥"
          mode="password"
          rules={[{ required: true, message: '请输入API密钥' }]}
        />
        <Form.InputNumber
          field="wechat_min_topup"
          label="最小充值金额"
          min={1}
        />
        <Button
          type="primary"
          htmlType="submit"
          loading={loading}
          block
        >
          保存配置
        </Button>
      </Form>

      <div className="mt-6 p-4 bg-gray-100 rounded">
        <h3 className="font-bold mb-2">配置说明：</h3>
        <ol className="list-decimal list-inside space-y-1 text-sm">
          <li>登录微信支付商户平台获取商户号和API密钥</li>
          <li>在商户平台配置支付回调地址：{'{服务器地址}/api/wechat/notify'}</li>
          <li>确保服务器外网可访问</li>
        </ol>
      </div>
    </div>
  );
}
```

#### `web/src/pages/Setting/Payment/SettingsPaymentGatewayAlipay.jsx`

类似微信支付配置页面，字段改为支付宝相关配置。

---

### 4.2 充值页面修改

在 `web/src/components/topup/index.jsx` 中添加支付方式：

```jsx
// 在 getTopupInfo 中处理新的支付方式
const enableWechatTopUp = data.enable_wechat_topup || false;
const enableAlipayTopUp = data.enable_alipay_topup || false;

setEnableWechatTopUp(enableWechatTopUp);
setEnableAlipayTopUp(enableAlipayTopUp);

// 添加支付方法到 payMethods
if (enableWechatTopUp) {
  payMethods.push({
    name: "微信支付",
    type: "wechat",
    color: "rgba(var(--semi-green-5), 1)",
    min_topup: data.wechat_min_topup
  });
}

if (enableAlipayTopUp) {
  payMethods.push({
    name: "支付宝",
    type: "alipay",
    color: "rgba(var(--semi-blue-5), 1)",
    min_topup: data.alipay_min_topup
  });
}

// 在 preTopUp 中添加处理逻辑
const preTopUp = async (payment) => {
    if (payment === 'wechat' && !enableWechatTopUp) {
        showError(t('管理员未开启微信支付！'));
        return;
    }
    if (payment === 'alipay' && !enableAlipayTopUp) {
        showError(t('管理员未开启支付宝支付！'));
        return;
    }

    setPayWay(payment);
    setPaymentLoading(true);

    try {
        if (payment === 'wechat' || payment === 'alipay') {
            // 获取金额
            await getAmount();
        }
        setOpen(true);
    } finally {
        setPaymentLoading(false);
    }
};

// 在 onlineTopUp 中添加微信/支付宝支付处理
const onlineTopUp = async () => {
    setConfirmLoading(true);
    try {
        let res;
        if (payWay === 'wechat') {
            res = await API.post('/api/user/wechat/pay', {
                amount: parseInt(topUpCount),
                payment_method: 'wechat',
            });
        } else if (payWay === 'alipay') {
            res = await API.post('/api/user/alipay/pay', {
                amount: parseInt(topUpCount),
                payment_method: 'alipay',
            });
        }

        if (res && res.data.message === 'success') {
            // 显示二维码支付
            window.open(res.data.data.pay_url, '_blank');
        }
    } catch (err) {
        showError(t('支付请求失败'));
    } finally {
        setConfirmLoading(false);
        setOpen(false);
    }
};
```

---

## 五、依赖安装

### Go 依赖

```bash
# 支付宝SDK
go get github.com/smartwalle/alipay/v3
```

微信支付使用原生实现，无需额外SDK。

---

## 六、数据库变更

无需修改数据库表结构，复用现有 `topup` 表，`payment_method` 字段新增值：
- `wechat` - 微信支付
- `alipay` - 支付宝

---

## 七、测试流程

### 7.1 微信支付测试
1. 配置商户号和API密钥
2. 用户发起充值，后端返回二维码链接
3. 使用微信扫码支付
4. 支付成功后，微信回调通知
5. 验证用户配额增加、订单状态更新

### 7.2 支付宝测试
1. 配置AppID和密钥
2. 用户发起充值，后端返回二维码链接
3. 使用支付宝扫码支付
4. 支付成功后，支付宝回调通知
5. 验证用户配额增加、订单状态更新

---

## 八、注意事项

1. **回调地址配置**：
   - 微信支付：在商户平台配置回调地址
   - 支付宝：在应用配置中设置回调地址

2. **金额单位**：
   - 微信支付使用"分"为单位
   - 支付宝使用"元"为单位（字符串格式，保留两位小数）

3. **签名验证**：
   - 微信支付：MD5签名 + API密钥
   - 支付宝：RSA2签名 + 公私钥

4. **订单幂等性**：
   - 使用订单锁机制防止重复处理
   - 回调时检查订单状态

5. **日志记录**：
   - 记录所有支付请求和回调
   - 便于排查问题

---

## 九、后续优化

1. 添加支付超时处理（定时任务关闭超时未支付订单）
2. 支持退款功能
3. 添加支付统计和对账功能
4. H5支付支持（移动端直接唤起微信/支付宝APP）
5. 小程序支付支持

---

## 十、常见问题

### Q1: 微信回调一直失败？
A: 检查回调地址是否配置正确、服务器是否可外网访问、签名验证逻辑

### Q2: 支付宝签名验证失败？
A: 检查应用私钥和支付宝公钥是否正确配置、格式是否正确

### Q3: 订单金额不一致？
A: 检查分组倍率、折扣配置、金额单位转换逻辑

### Q4: 用户配额没有增加？
A: 检查回调处理逻辑、数据库事务是否成功、日志记录
