# 微信登录二维码显示问题 - 解决方案

## 问题现象

用户在使用微信登录时，微信公众号二维码无法显示，显示为损坏的图片图标。

## 问题原因

系统配置中的 `WeChatAccountQRCodeImageURL` 字段未配置或为空字符串，导致前端无法加载二维码图片。

## 技术细节

### 代码路径

**前端显示** (web/src/components/auth/LoginForm.jsx:813):
```jsx
<img src={status.wechat_qrcode} alt='微信二维码' className='mb-4' />
```

**后端API** (controller/misc.go:65):
```go
"wechat_qrcode": common.WeChatAccountQRCodeImageURL,
```

**配置初始化** (model/option.go:112):
```go
common.OptionMap["WeChatAccountQRCodeImageURL"] = ""  // 默认为空
```

**配置更新** (model/option.go:407-408):
```go
case "WeChatAccountQRCodeImageURL":
    common.WeChatAccountQRCodeImageURL = value
```

### 数据流程

```
初始化
  ↓
WeChatAccountQRCodeImageURL = ""
  ↓
GET /api/status
  ↓
返回 wechat_qrcode: ""
  ↓
<img src="" />  ← 空src导致显示损坏图标
```

---

## 解决方案

### 方法一：管理后台配置（推荐）

#### 1. 登录管理后台

使用管理员账号登录系统。

#### 2. 进入系统设置

导航路径：**设置 → 系统设置**

#### 3. 配置微信公众号二维码

找到 **"微信公众号二维码图片链接"** (WeChatAccountQRCodeImageURL) 字段，填入二维码图片的公网URL。

**示例配置**:
```
https://your-domain.com/images/wechat_official_qrcode.png
```

**要求**:
- ✅ 必须是公网可访问的URL
- ✅ 推荐使用HTTPS链接
- ✅ 支持的图片格式：PNG、JPG、GIF
- ✅ 建议尺寸：200x200 或更大（系统会自动缩放）

#### 4. 保存配置

点击"保存"按钮。

#### 5. 验证

刷新登录页面，点击"微信登录"按钮，确认二维码正常显示。

---

### 方法二：数据库直接配置

适用于批量部署或自动化场景。

#### PostgreSQL

```sql
-- 插入或更新配置
INSERT INTO options (key, value)
VALUES ('WeChatAccountQRCodeImageURL', 'https://your-domain.com/images/wechat_qrcode.png')
ON CONFLICT (key)
DO UPDATE SET value = EXCLUDED.value;

-- 验证配置
SELECT * FROM options WHERE key = 'WeChatAccountQRCodeImageURL';
```

#### MySQL

```sql
-- 插入或更新配置
INSERT INTO options (`key`, `value`)
VALUES ('WeChatAccountQRCodeImageURL', 'https://your-domain.com/images/wechat_qrcode.png')
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- 验证配置
SELECT * FROM options WHERE `key` = 'WeChatAccountQRCodeImageURL';
```

#### SQLite

```sql
-- 插入或更新配置
INSERT OR REPLACE INTO options (key, value)
VALUES ('WeChatAccountQRCodeImageURL', 'https://your-domain.com/images/wechat_qrcode.png');

-- 验证配置
SELECT * FROM options WHERE key = 'WeChatAccountQRCodeImageURL';
```

#### Docker环境执行

**PostgreSQL**:
```bash
docker exec -it postgres psql -U root -d new-api -c \
  "INSERT INTO options (key, value) VALUES ('WeChatAccountQRCodeImageURL', 'https://your-domain.com/images/wechat_qrcode.png') ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value;"
```

**MySQL**:
```bash
docker exec -it mysql mysql -uroot -p123456 new-api -e \
  "INSERT INTO options (\`key\`, \`value\`) VALUES ('WeChatAccountQRCodeImageURL', 'https://your-domain.com/images/wechat_qrcode.png') ON DUPLICATE KEY UPDATE \`value\` = VALUES(\`value\`);"
```

---

### 方法三：环境变量配置（临时测试）

在 `.env` 文件或 `docker-compose.yml` 中设置（**注意**：仅适用于首次启动，后续会被数据库配置覆盖）：

```bash
WECHAT_ACCOUNT_QRCODE_IMAGE_URL=https://your-domain.com/images/wechat_qrcode.png
```

**重启服务**:
```bash
docker-compose restart new-api
```

---

## 二维码图片准备

### 1. 获取微信公众号二维码

1. 登录 [微信公众平台](https://mp.weixin.qq.com/)
2. 进入 **设置与开发 → 公众号设置 → 账号详情**
3. 找到 **"公众号二维码"**，点击下载

### 2. 上传到服务器

**方法A: 使用项目静态资源目录**

```bash
# 创建静态图片目录（如果不存在）
mkdir -p web/public/images

# 上传二维码图片
cp /path/to/wechat_qrcode.png web/public/images/

# 重新构建前端
cd web
bun run build

# 二维码访问URL
https://your-domain.com/images/wechat_qrcode.png
```

**方法B: 使用第三方图床**

可以使用以下图床服务上传二维码：
- 阿里云OSS
- 腾讯云COS
- 七牛云
- 又拍云
- GitHub Pages (公开仓库)

示例（GitHub）:
1. 创建一个公开的GitHub仓库
2. 上传二维码到仓库
3. 使用raw链接：`https://raw.githubusercontent.com/username/repo/main/wechat_qrcode.png`

---

## 验证配置

### 1. 检查后端API

```bash
# 调用状态API查看配置
curl http://localhost:3000/api/status | jq .wechat_qrcode

# 预期输出（配置成功）:
"https://your-domain.com/images/wechat_qrcode.png"

# 预期输出（配置失败）:
""
```

### 2. 检查前端控制台

打开浏览器开发者工具（F12）→ Network标签，点击微信登录按钮：
- 查看 `status` API 响应中的 `wechat_qrcode` 字段
- 查看二维码图片请求是否成功（状态码应为200）

### 3. 检查数据库

```sql
SELECT * FROM options WHERE key = 'WeChatAccountQRCodeImageURL';
```

预期结果：
```
key                            | value
-------------------------------|-------------------------------------------
WeChatAccountQRCodeImageURL    | https://your-domain.com/images/wechat_qrcode.png
```

---

## 常见问题

### Q1: 配置后二维码仍不显示？

**可能原因**:
1. **缓存问题** - 浏览器缓存了旧的状态数据
   - **解决**: 硬刷新页面（Ctrl+Shift+R 或 Cmd+Shift+R）
   - **解决**: 清除浏览器缓存后重新登录

2. **图片URL无法访问** - 防火墙、CORS或权限问题
   - **解决**: 直接在浏览器打开URL，确认可访问
   - **解决**: 检查服务器防火墙和Nginx配置

3. **HTTPS混合内容错误** - HTTPS站点加载HTTP图片被浏览器阻止
   - **解决**: 确保图片URL使用HTTPS

### Q2: 如何快速测试图片URL是否有效？

```bash
# 使用curl测试
curl -I https://your-domain.com/images/wechat_qrcode.png

# 预期输出（成功）:
HTTP/2 200
content-type: image/png
content-length: 12345
...
```

### Q3: Docker环境下静态文件如何部署？

**方法A: 挂载静态资源目录**

修改 `docker-compose.yml`:
```yaml
services:
  new-api:
    volumes:
      - ./web/public:/app/web/public
```

**方法B: 重新构建镜像**

```bash
# 将二维码放入web/public/images/目录
cp wechat_qrcode.png web/public/images/

# 重新构建
docker-compose build new-api
docker-compose up -d new-api
```

### Q4: 配置后多久生效？

- **立即生效**: 配置保存后即刻写入数据库
- **缓存同步**: 默认60秒同步一次（可通过 `SYNC_FREQUENCY` 环境变量调整）
- **建议**: 配置后等待1-2分钟，或重启服务确保生效

---

## 最佳实践

1. **使用HTTPS图片链接** - 避免混合内容警告
2. **优化图片大小** - 建议200x200 PNG格式，文件大小<100KB
3. **使用CDN** - 提高加载速度和可用性
4. **定期检查** - 确保图片URL长期有效（避免外链失效）
5. **备份二维码** - 本地保留备份，防止丢失

---

## 相关文件

| 文件路径 | 作用 |
|---------|------|
| `web/src/components/auth/LoginForm.jsx` | 前端登录表单，显示二维码 |
| `web/src/components/settings/personal/modals/WeChatBindModal.jsx` | 微信绑定弹窗，显示二维码 |
| `web/src/components/settings/SystemSetting.jsx` | 系统设置页面，配置二维码URL |
| `controller/misc.go` | 后端状态API，返回配置给前端 |
| `controller/wechat.go` | 微信认证逻辑 |
| `model/option.go` | 配置管理模块 |
| `common/constants.go` | 全局常量定义 |

---

## 总结

微信登录二维码无法显示是因为 **`WeChatAccountQRCodeImageURL` 配置为空**。

**快速解决步骤**:
1. 登录管理后台 → 系统设置
2. 找到"微信公众号二维码图片链接"字段
3. 填入二维码图片的公网URL
4. 保存并刷新登录页面

配置完成后，微信登录二维码将正常显示。
