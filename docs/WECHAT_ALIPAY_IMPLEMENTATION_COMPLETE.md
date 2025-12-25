# 微信支付和支付宝支付 - 完整实现报告

## ✅ 已完成的工作

### 1. 后端实现（已全部完成）

#### 配置模块
- ✅ `setting/payment_wechat.go` - 微信支付配置
  - WechatAppId（可选，用于 JSAPI 支付）
  - WechatMchId（商户号）
  - WechatApiV2Key（API v2 密钥）
  - WechatMinTopUp（最小充值金额）

- ✅ `setting/payment_alipay.go` - 支付宝配置
  - AlipayAppId（应用 AppID）
  - AlipayPrivateKey（应用私钥）
  - AlipayPublicKey（支付宝公钥）
  - AlipayMinTopUp（最小充值金额）

#### 控制器
- ✅ `controller/topup_wechat.go` - 微信支付控制器（445行）
  - `RequestWechatPay()` - 处理支付请求
  - `wechatUnifiedOrder()` - 调用微信统一下单 API
  - `wechatSign()` - MD5 签名生成
  - `verifyWechatSign()` - 签名验证
  - `WechatNotifyHandler()` - 支付回调处理
  - `getWechatPayMoney()` - 金额计算（支持分组倍率和折扣）
  - `RequestWechatAmount()` - 获取支付金额（前端显示）
  - 支持 Native 扫码支付（二维码）
  - 完整的订单幂等性控制

- ✅ `controller/topup_alipay.go` - 支付宝控制器（327行）
  - `InitAlipayClient()` - 初始化支付宝 SDK 客户端
  - `RequestAlipayPay()` - 处理支付请求
  - `alipayQRCodePay()` - 调用支付宝当面付（扫码支付）
  - `AlipayNotifyHandler()` - 支付回调处理
  - `getAlipayPayMoney()` - 金额计算
  - `RequestAlipayAmount()` - 获取支付金额（前端显示）
  - 使用 `github.com/smartwalle/alipay/v3` SDK
  - RSA2 签名验证
  - 完整的订单幂等性控制

#### 路由配置
- ✅ 更新 `router/api-router.go`
  - 添加微信支付回调路由: `POST /api/wechat/notify`
  - 添加支付宝回调路由: `POST /api/alipay/notify`
  - 添加用户支付请求路由:
    - `POST /api/user/wechat/pay`
    - `POST /api/user/wechat/amount`
    - `POST /api/user/alipay/pay`
    - `POST /api/user/alipay/amount`

#### 充值信息接口
- ✅ 更新 `controller/topup.go` 中的 `GetTopUpInfo` 函数
  - 添加微信支付配置检测（第51-70行）
  - 添加支付宝配置检测（第72-91行）
  - 返回前端所需的支付方式列表和配置信息

#### 依赖管理
- ✅ 安装支付宝 SDK: `github.com/smartwalle/alipay/v3 v3.2.28`
- ✅ 代码编译成功，无错误

---

### 2. 前端实现（已全部完成）

#### 配置页面
- ✅ `web/src/pages/Setting/Payment/SettingsPaymentGatewayWechat.jsx`
  - 微信支付配置表单
  - 包含字段：AppID、商户号、API v2 密钥、最小充值金额
  - 显示回调地址提示
  - 表单验证和错误处理

- ✅ `web/src/pages/Setting/Payment/SettingsPaymentGatewayAlipay.jsx`
  - 支付宝配置表单
  - 包含字段：AppID、应用私钥、支付宝公钥、最小充值金额
  - 多行文本输入框用于密钥配置
  - 显示回调地址提示
  - 表单验证和错误处理

#### 集成配置页面
- ✅ 更新 `web/src/components/settings/PaymentSetting.jsx`
  - 导入微信和支付宝配置组件
  - 添加状态管理（WechatAppId, WechatMchId, WechatApiV2Key, WechatMinTopUp）
  - 添加状态管理（AlipayAppId, AlipayPrivateKey, AlipayPublicKey, AlipayMinTopUp）
  - 在 getOptions 中解析配置（支持数值类型转换）
  - 渲染微信和支付宝配置卡片

#### 充值组件
- ✅ 更新 `web/src/components/topup/index.jsx`
  - 添加 `getWechatAmount()` 函数 - 获取微信支付金额
  - 添加 `getAlipayAmount()` 函数 - 获取支付宝支付金额
  - 更新 `preTopUp()` 函数 - 支持微信和支付宝支付前的金额计算
  - 更新 `onlineTopUp()` 函数 - 完整支付流程：
    - 调用 `/api/user/wechat/pay` 或 `/api/user/alipay/pay`
    - 获取二维码链接（data.pay_url）
    - 打开新窗口显示支付二维码
    - 显示成功提示："支付二维码已打开，请扫码支付"

#### 代码质量
- ✅ 运行 `bun run lint:fix` 进行代码格式化
- ✅ 遵循项目代码规范

---

## 🔧 配置指南

### 后端配置

#### 方式1: 环境变量配置（可选）

在 `.env` 文件或系统环境变量中配置：

```bash
# 微信支付配置
WECHAT_APP_ID=wx1234567890abcdef     # 可选
WECHAT_MCH_ID=1234567890             # 必填
WECHAT_API_V2_KEY=your_api_v2_key    # 必填
WECHAT_MIN_TOPUP=1

# 支付宝配置
ALIPAY_APP_ID=2021001234567890
ALIPAY_PRIVATE_KEY=MIIEvQIBADANBgkq...  # 应用私钥（PKCS8格式）
ALIPAY_PUBLIC_KEY=MIIBIjANBgkqhki...    # 支付宝公钥
ALIPAY_MIN_TOPUP=1
```

#### 方式2: 管理后台配置（推荐）

1. 登录管理后台
2. 进入 "设置" → "支付设置"
3. 找到 "微信支付设置" 部分，填写：
   - 商户号
   - API v2 密钥
   - 最小充值金额
4. 找到 "支付宝设置" 部分，填写：
   - 应用 AppID
   - 应用私钥
   - 支付宝公钥
   - 最小充值金额
5. 点击"更新"按钮保存配置

### 微信支付商户平台配置

1. 登录 [微信支付商户平台](https://pay.weixin.qq.com/)
2. 产品中心 → 开发配置 → API安全
3. 获取商户号和 API v2 密钥
4. 配置支付回调地址：
   ```
   https://yourdomain.com/api/wechat/notify
   ```
5. 确保服务器外网可访问

### 支付宝开放平台配置

1. 登录 [支付宝开放平台](https://open.alipay.com/)
2. 创建应用并获取 AppID
3. 配置应用私钥和支付宝公钥：
   - 生成应用私钥（PKCS8 格式）
   - 上传应用公钥到支付宝平台
   - 下载支付宝公钥
4. 添加功能：当面付（扫码支付）
5. 配置支付回调地址：
   ```
   https://yourdomain.com/api/alipay/notify
   ```
6. 建议先在沙箱环境测试

---

## 📊 完整数据流程

### 微信支付流程

```
1. 用户选择微信支付 → 输入充值金额
   ↓
2. 前端调用 POST /api/user/wechat/amount（获取实际支付金额）
   ↓
3. 用户确认 → 前端调用 POST /api/user/wechat/pay
   ↓
4. 后端创建订单（status=pending）
   ↓
5. 调用微信统一下单 API
   - 请求参数：商户号、订单号、金额（分）、回调地址等
   - MD5 签名验证
   ↓
6. 返回二维码链接（code_url）
   ↓
7. 前端打开新窗口显示二维码
   ↓
8. 用户扫码支付
   ↓
9. 微信回调 POST /api/wechat/notify
   - 验证签名
   - 订单加锁（防止并发）
   - 检查订单状态（幂等性）
   ↓
10. 增加用户配额
   ↓
11. 更新订单状态（status=success）
   ↓
12. 返回 SUCCESS 给微信
```

### 支付宝支付流程

```
1. 用户选择支付宝 → 输入充值金额
   ↓
2. 前端调用 POST /api/user/alipay/amount（获取实际支付金额）
   ↓
3. 用户确认 → 前端调用 POST /api/user/alipay/pay
   ↓
4. 后端创建订单（status=pending）
   ↓
5. 调用支付宝预下单 API（当面付）
   - 请求参数：AppID、订单号、金额（元）、回调地址等
   - RSA2 签名
   ↓
6. 返回二维码字符串（qr_code）
   ↓
7. 前端打开新窗口显示二维码
   ↓
8. 用户扫码支付
   ↓
9. 支付宝回调 POST /api/alipay/notify
   - RSA2 验签
   - 订单加锁（防止并发）
   - 检查订单状态（幂等性）
   ↓
10. 增加用户配额
   ↓
11. 更新订单状态（status=success）
   ↓
12. 返回 success 给支付宝
```

---

## 🔐 安全机制

### 1. 签名验证
- **微信支付**: MD5 签名 + API v2 密钥
- **支付宝**: RSA2 签名 + 公私钥对

### 2. 订单幂等性
- 使用订单锁（sync.Map）防止并发处理
- 回调时检查订单状态，避免重复充值
- 代码位置：`controller/topup.go` 第251-276行

```go
LockOrder(tradeNo)
defer UnlockOrder(tradeNo)

if topUp.Status != common.TopUpStatusPending {
    // 订单已处理，直接返回成功
    return
}
```

### 3. 金额验证
- 后端计算实际支付金额（不信任前端）
- 支持分组倍率（TopupGroupRatio）
- 支持折扣配置（AmountDiscount）

### 4. 配额计算
- 支持 Token 和 USD 两种显示方式
- 准确的 QuotaPerUnit 转换
- 代码位置：
  - 微信：`controller/topup_wechat.go` 第324-327行
  - 支付宝：`controller/topup_alipay.go` 第219-221行

---

## 🧪 测试步骤

### 准备工作

1. **启动服务器**
   ```bash
   go run main.go
   ```

2. **确保前端已构建**
   ```bash
   cd web
   bun run build
   ```

3. **配置回调地址**
   - 确保服务器外网可访问（使用公网 IP 或域名）
   - 或使用内网穿透工具（如 ngrok）进行本地测试

### 微信支付测试

1. **配置微信支付**
   - 登录管理后台 → 设置 → 支付设置
   - 填写商户号和 API v2 密钥
   - 保存配置

2. **发起充值**
   - 登录用户账号
   - 进入充值页面
   - 选择"微信支付"
   - 输入充值金额（≥ 最小充值金额）
   - 点击确认

3. **扫码支付**
   - 新窗口打开，显示二维码
   - 使用微信"扫一扫"扫描二维码
   - 完成支付

4. **验证结果**
   - 检查用户配额是否增加
   - 查看充值历史记录
   - 检查订单状态（应为 success）
   - 查看服务器日志：
     ```
     微信支付成功 - 订单号: xxx, 微信流水号: xxx, 金额: x.xx, 配额: xxxx
     ```

### 支付宝测试

1. **配置支付宝**
   - 登录管理后台 → 设置 → 支付设置
   - 填写 AppID、应用私钥、支付宝公钥
   - 保存配置

2. **发起充值**
   - 登录用户账号
   - 进入充值页面
   - 选择"支付宝"
   - 输入充值金额（≥ 最小充值金额）
   - 点击确认

3. **扫码支付**
   - 新窗口打开，显示二维码
   - 使用支付宝"扫一扫"扫描二维码
   - 完成支付

4. **验证结果**
   - 检查用户配额是否增加
   - 查看充值历史记录
   - 检查订单状态（应为 success）
   - 查看服务器日志：
     ```
     支付宝支付成功 - 订单号: xxx, 支付宝流水号: xxx, 金额: x.xx, 配额: xxxx
     ```

### 沙箱环境测试

#### 支付宝沙箱

1. 登录支付宝开放平台
2. 进入"开发者中心" → "研发服务" → "沙箱环境"
3. 获取沙箱 AppID 和密钥
4. 修改 `controller/topup_alipay.go` 第33行：
   ```go
   alipayClient, err = alipay.New(
       setting.AlipayAppId,
       setting.AlipayPrivateKey,
       true, // 改为 true 使用沙箱环境
   )
   ```
5. 使用沙箱账号进行测试

---

## 🐛 故障排查

### 常见问题

#### 1. 微信回调失败

**现象**: 支付成功但配额未增加

**排查步骤**:
```bash
# 查看服务器日志
docker-compose logs -f new-api
# 或
tail -f ./logs/error.log
```

**可能原因**:
- 回调地址配置错误
- 服务器外网不可访问
- 签名验证失败（API 密钥错误）
- 防火墙拦截

**解决方案**:
1. 确认回调地址：`https://yourdomain.com/api/wechat/notify`
2. 测试服务器外网可访问性：`curl https://yourdomain.com/api/status`
3. 检查 API 密钥是否正确
4. 检查日志中的详细错误信息

#### 2. 支付宝签名验证失败

**现象**: 回调返回 "fail"，日志显示 "签名验证失败"

**可能原因**:
- 应用私钥格式错误（应为 PKCS8）
- 支付宝公钥配置错误
- AppID 不匹配

**解决方案**:
1. 检查应用私钥格式：
   ```bash
   # PKCS8 格式示例（正确）
   MIIEvQIBADANBgkqhki...

   # PKCS1 格式（错误，需要转换）
   MIICXAIBAAKBgQC...
   ```
2. 重新生成密钥对并上传到支付宝平台
3. 确认 AppID 与应用一致

#### 3. 金额计算错误

**现象**: 实际扣费金额与预期不符

**排查步骤**:
```go
// 检查配置
operation_setting.Price           // 基础价格
common.GetTopupGroupRatio(group) // 分组倍率
operation_setting.GetPaymentSetting().AmountDiscount  // 折扣配置
common.QuotaPerUnit              // 配额单位
```

**计算公式**:
```
quota = amount × Price × TopupGroupRatio × Discount
payMoney = quota × QuotaPerUnit
```

#### 4. 订单状态未更新

**现象**: 订单一直处于 pending 状态

**可能原因**:
- 回调未触发
- 回调处理过程中发生错误
- 数据库事务失败

**排查步骤**:
1. 检查回调是否触发（查看日志）
2. 查看错误日志中的详细信息
3. 检查数据库事务是否成功
4. 手动查询订单：
   ```sql
   SELECT * FROM topups WHERE trade_no = 'xxx';
   ```

---

## 📝 日志记录

所有支付操作都会记录详细日志，方便问题排查：

### 微信支付日志

```go
// 支付成功
log.Printf("微信支付成功 - 订单号: %s, 微信流水号: %s, 金额: %.2f, 配额: %d",
    notify.OutTradeNo, notify.TransactionId, topUp.Money, quotaToAdd)

// 签名验证失败
log.Printf("微信回调签名验证失败")

// 订单不存在
log.Printf("微信回调: 订单不存在 - %s", notify.OutTradeNo)

// 订单已处理
log.Printf("微信回调: 订单已处理 - %s, 状态: %s", notify.OutTradeNo, topUp.Status)
```

### 支付宝日志

```go
// 支付成功
log.Printf("支付宝支付成功 - 订单号: %s, 支付宝流水号: %s, 金额: %.2f, 配额: %d",
    outTradeNo, tradeNo, topUp.Money, quotaToAdd)

// 回调状态异常
log.Printf("支付宝回调状态异常: %s - %s", outTradeNo, tradeStatus)

// 订单不存在
log.Printf("支付宝回调: 订单不存在 - %s", outTradeNo)

// 订单已处理
log.Printf("支付宝回调: 订单已处理 - %s, 状态: %s", outTradeNo, topUp.Status)
```

### 查看日志

```bash
# Docker 部署
docker-compose logs -f new-api

# 本地运行
tail -f ./logs/error.log

# 查看最近 100 行
docker-compose logs --tail 100 new-api
```

---

## 💡 注意事项

### 开发环境

1. **使用沙箱测试**
   - 微信支付：暂无官方沙箱环境，建议用小额测试
   - 支付宝：使用开放平台提供的沙箱环境

2. **内网穿透**
   - 本地开发时使用 ngrok 等工具暴露回调地址
   ```bash
   ngrok http 3000
   # 将生成的 https 地址配置到支付平台
   ```

### 生产环境

1. **HTTPS 必须**
   - 回调地址必须使用 HTTPS
   - 配置 SSL 证书

2. **密钥安全**
   - 不要在代码中硬编码密钥
   - 使用环境变量或密钥管理服务
   - 定期更换密钥

3. **监控和告警**
   - 监控支付成功率
   - 设置订单超时告警
   - 定期对账

4. **定期对账**
   - 下载微信/支付宝对账单
   - 与系统订单记录比对
   - 处理异常订单

---

## ✨ 功能特性

- ✅ 支持微信Native扫码支付
- ✅ 支持支付宝当面付扫码
- ✅ 自动计算实际支付金额（支持分组倍率和折扣）
- ✅ 完整的签名验证机制
- ✅ 订单幂等性控制（防止重复充值）
- ✅ 详细的日志记录
- ✅ 前后端完整实现
- ✅ 管理后台配置界面
- ✅ 支持最小充值金额设置

---

## 📚 相关文档

- [微信支付官方文档](https://pay.weixin.qq.com/wiki/doc/api/index.html)
- [支付宝开放平台文档](https://opendocs.alipay.com/)
- [支付宝 Go SDK](https://github.com/smartwalle/alipay)

---

## 🎉 完成！

微信支付和支付宝支付功能已全部实现完成，包括：

1. ✅ 后端完整实现（配置、控制器、路由、回调处理）
2. ✅ 前端完整实现（配置页面、充值组件）
3. ✅ 代码编译通过
4. ✅ 遵循项目代码规范
5. ✅ 完整的文档和测试指南

**下一步**: 配置支付平台并进行测试！

如有问题，请查看日志或参考故障排查部分。
