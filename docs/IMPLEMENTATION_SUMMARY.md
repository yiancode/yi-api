# 微信支付和支付宝支付 - 实现总结

## ✅ 已完成的工作

### 1. 后端代码实现（已完成）

#### 配置模块
- ✅ `setting/payment_wechat.go` - 微信支付配置
- ✅ `setting/payment_alipay.go` - 支付宝配置

#### 控制器
- ✅ `controller/topup_wechat.go` - 微信支付控制器（完整实现）
  - 支付请求处理
  - 微信统一下单API调用
  - 签名生成和验证
  - 支付回调处理
  - 订单幂等性控制

- ✅ `controller/topup_alipay.go` - 支付宝支付控制器（完整实现）
  - 支付请求处理
  - 支付宝当面付API调用
  - 支付回调处理
  - 订单幂等性控制
  - 已修复 SDK 兼容性问题

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
  - 添加微信支付配置检测
  - 添加支付宝配置检测
  - 返回前端所需的支付方式列表

#### 依赖管理
- ✅ 已安装 Go 依赖: `github.com/smartwalle/alipay/v3 v3.2.28`
- ✅ 运行 `go mod tidy` 成功
- ✅ 代码编译通过，无错误

---

### 2. 前端代码实现（已完成）

#### 配置页面
- ✅ `web/src/pages/Setting/Payment/SettingsPaymentGatewayWechat.jsx`
  - 微信支付配置表单
  - 字段：AppID（可选）、商户号、API v2 密钥、最小充值金额
  - 显示回调地址和配置说明
  - 完整的表单验证和错误处理

- ✅ `web/src/pages/Setting/Payment/SettingsPaymentGatewayAlipay.jsx`
  - 支付宝配置表单
  - 字段：AppID、应用私钥、支付宝公钥、最小充值金额
  - 多行文本框输入密钥
  - 显示回调地址和配置说明
  - 完整的表单验证和错误处理

#### 集成到主配置页面
- ✅ 更新 `web/src/components/settings/PaymentSetting.jsx`
  - 导入微信和支付宝配置组件
  - 添加状态管理（8个新配置项）
  - 在 `getOptions` 中解析配置（支持数值类型）
  - 渲染微信和支付宝配置卡片

#### 充值组件
- ✅ 更新 `web/src/components/topup/index.jsx`
  - 添加 `getWechatAmount()` 函数
  - 添加 `getAlipayAmount()` 函数
  - 更新 `preTopUp()` 函数支持微信和支付宝
  - 更新 `onlineTopUp()` 函数完整支付流程：
    - 调用对应的支付 API
    - 打开二维码链接
    - 显示成功提示

#### 代码质量
- ✅ 运行 `bun run lint:fix` 进行代码格式化
- ✅ 遵循项目代码规范

---

## 🎉 实现完成

所有功能已完整实现，包括：
- ✅ 后端完整实现（配置、控制器、路由）
- ✅ 前端完整实现（配置页面、充值组件）
- ✅ 代码编译通过
- ✅ 代码格式化完成

---

```bash
cd /Volumes/SSD/ssd-code/github/yi-api

# 安装支付宝 SDK
go get github.com/smartwalle/alipay/v3

# 更新依赖
go mod tidy
```

### 步骤 2: 配置环境变量（可选）

在 `.env` 文件或系统环境变量中配置支付参数（也可以通过管理后台配置）：

```bash
# 微信支付配置
WECHAT_MCH_ID=你的商户号
WECHAT_API_V2_KEY=你的API密钥
WECHAT_MIN_TOPUP=1

# 支付宝配置
ALIPAY_APP_ID=你的AppID
ALIPAY_PRIVATE_KEY=你的应用私钥
ALIPAY_PUBLIC_KEY=支付宝公钥
ALIPAY_MIN_TOPUP=1
```

### 步骤 3: 前端配置界面（参考实现）

创建 `web/src/pages/Setting/Payment/SettingsPaymentGatewayWechat.jsx`：

```jsx
import React, { useState, useEffect } from 'react';
import { Form, Input, Button, Toast, InputNumber } from '@douyinfe/semi-ui';
import { API } from '../../../helpers';

export default function SettingsPaymentGatewayWechat() {
  const [loading, setLoading] = useState(false);

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
      <Form onSubmit={handleSubmit} style={{ maxWidth: 600 }}>
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
          initValue={1}
          min={1}
        />
        <Button type="primary" htmlType="submit" loading={loading} block>
          保存配置
        </Button>
      </Form>

      <div className="mt-6 p-4 bg-gray-100 rounded">
        <h3 className="font-bold mb-2">配置说明：</h3>
        <ol className="list-decimal list-inside space-y-1 text-sm">
          <li>登录微信支付商户平台获取商户号和API密钥</li>
          <li>配置回调地址：{'{服务器地址}/api/wechat/notify'}</li>
          <li>确保服务器外网可访问</li>
        </ol>
      </div>
    </div>
  );
}
```

创建类似的 `SettingsPaymentGatewayAlipay.jsx` 用于支付宝配置。

### 步骤 4: 更新前端充值页面

在 `web/src/components/topup/index.jsx` 中的 `onlineTopUp` 函数添加微信和支付宝处理：

```jsx
const onlineTopUp = async () => {
    setConfirmLoading(true);
    try {
        let res;

        if (payWay === 'wechat') {
            // 微信支付
            res = await API.post('/api/user/wechat/pay', {
                amount: parseInt(topUpCount),
                payment_method: 'wechat',
            });
        } else if (payWay === 'alipay') {
            // 支付宝支付
            res = await API.post('/api/user/alipay/pay', {
                amount: parseInt(topUpCount),
                payment_method: 'alipay',
            });
        } else if (payWay === 'stripe') {
            // 原有Stripe逻辑
            res = await API.post('/api/user/stripe/pay', {
                amount: parseInt(topUpCount),
                payment_method: 'stripe',
            });
        } else {
            // 原有易支付逻辑
            res = await API.post('/api/user/pay', {
                amount: parseInt(topUpCount),
                payment_method: payWay,
            });
        }

        if (res && res.data.message === 'success') {
            if (payWay === 'wechat' || payWay === 'alipay' || payWay === 'stripe') {
                // 打开支付链接（二维码）
                window.open(res.data.data.pay_url, '_blank');
                showSuccess(t('支付二维码已打开，请扫码支付'));
            } else {
                // 原有表单提交逻辑...
            }
        } else {
            showError(res.data.data || '支付请求失败');
        }
    } catch (err) {
        console.error(err);
        showError(t('支付请求失败'));
    } finally {
        setConfirmLoading(false);
        setOpen(false);
    }
};
```

### 步骤 5: 测试流程

#### 微信支付测试
1. 在管理后台配置微信商户号和API密钥
2. 用户发起充值，选择微信支付
3. 后端返回二维码链接
4. 使用微信扫码支付
5. 支付成功后，微信回调 `/api/wechat/notify`
6. 验证用户配额增加、订单状态更新

#### 支付宝测试
1. 在管理后台配置支付宝AppID和密钥
2. 用户发起充值，选择支付宝
3. 后端返回二维码链接
4. 使用支付宝扫码支付
5. 支付成功后，支付宝回调 `/api/alipay/notify`
6. 验证用户配额增加、订单状态更新

---

## 🔧 配置指南

### 微信支付商户平台配置

1. 登录 [微信支付商户平台](https://pay.weixin.qq.com/)
2. 产品中心 → 开发配置 → API安全 → 下载API证书
3. 获取商户号和API密钥
4. 配置支付回调地址：`https://yourdomain.com/api/wechat/notify`

### 支付宝开放平台配置

1. 登录 [支付宝开放平台](https://open.alipay.com/)
2. 创建应用并获取AppID
3. 配置应用私钥和支付宝公钥
4. 添加功能：当面付
5. 配置支付回调地址：`https://yourdomain.com/api/alipay/notify`

---

## 📊 数据流程

### 微信支付流程

```
用户发起充值
    ↓
POST /api/user/wechat/pay
    ↓
创建订单(pending状态)
    ↓
调用微信统一下单API
    ↓
返回二维码链接
    ↓
用户扫码支付
    ↓
微信回调 POST /api/wechat/notify
    ↓
验证签名 → 查询订单 → 订单加锁
    ↓
增加用户配额 → 更新订单(success状态)
    ↓
返回SUCCESS给微信
```

### 支付宝支付流程

```
用户发起充值
    ↓
POST /api/user/alipay/pay
    ↓
创建订单(pending状态)
    ↓
调用支付宝预下单API(当面付)
    ↓
返回二维码链接
    ↓
用户扫码支付
    ↓
支付宝回调 POST /api/alipay/notify
    ↓
验证签名 → 查询订单 → 订单加锁
    ↓
增加用户配额 → 更新订单(success状态)
    ↓
返回success给支付宝
```

---

## 🔐 安全机制

### 1. 签名验证
- 微信：MD5签名 + API密钥
- 支付宝：RSA2签名 + 公私钥对

### 2. 订单幂等性
- 使用订单锁（sync.Map）防止并发处理
- 回调时检查订单状态，避免重复充值

### 3. 金额验证
- 后端计算实际支付金额
- 支持分组倍率、折扣等

### 4. 配额计算
- 支持Token和USD两种显示方式
- 准确的QuotaPerUnit转换

---

## 🐛 故障排查

### 常见问题

**Q: 微信回调一直失败？**
- 检查回调地址配置
- 检查服务器外网可访问性
- 查看日志中的签名验证错误
- 确认API密钥配置正确

**Q: 支付宝签名验证失败？**
- 检查应用私钥格式（PKCS1/PKCS8）
- 检查支付宝公钥是否正确
- 确认AppID配置正确

**Q: 订单状态没有更新？**
- 检查回调是否触发
- 查看日志中的错误信息
- 检查数据库事务是否成功

**Q: 用户配额没有增加？**
- 检查IncreaseUserQuota函数调用
- 查看日志记录
- 检查QuotaPerUnit计算

---

## 📝 日志记录

所有支付操作都会记录详细日志：

```go
log.Printf("微信支付成功 - 订单号: %s, 微信流水号: %s, 金额: %.2f, 配额: %d", ...)
log.Printf("支付宝支付成功 - 订单号: %s, 支付宝流水号: %s, 金额: %.2f, 配额: %d", ...)
```

使用 `docker-compose logs -f new-api` 或 `tail -f ./logs/error.log` 查看日志。

---

## ✨ 完成！

现在你的 yi-api 项目已经支持微信支付和支付宝支付了！

**重要提醒**：
1. 测试环境先使用沙箱测试
2. 生产环境确保回调地址HTTPS
3. 保护好商户密钥和私钥文件
4. 定期查看支付日志和对账

如有问题，请查看详细实现方案文档：`docs/WECHAT_ALIPAY_PAYMENT_IMPLEMENTATION.md`
