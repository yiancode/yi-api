# Creem 支付合规性部署清单

本文档列出为满足 Creem 支付合规要求而进行的所有更改，以及部署步骤。

## ✅ 已完成的更改

### 1. 法律文档创建
- ✅ 创建了完整的隐私政策（Privacy Policy）
  - 文件位置：`docs/privacy_policy.md`
  - 包含：数据收集、使用、分享、安全、用户权利、GDPR 合规

- ✅ 创建了完整的服务条款（Terms of Service）
  - 文件位置：`docs/terms_of_service.md`
  - 包含：服务描述、可接受使用、计费条款、退款政策、责任免责声明、争议解决

### 2. 前端更新
- ✅ 修改了 Footer 组件（`web/src/components/layout/Footer.jsx`）
  - 添加了"隐私政策"链接：`/privacy`
  - 添加了"服务条款"链接：`/user-agreement`
  - 添加了"客户支持"邮箱链接：`mailto:your-email@example.com`

- ✅ 前端已构建（`bun run build` 成功）

### 3. 回复邮件准备
- ✅ 邮件草稿已创建：`docs/CREEM_REPLY_EMAIL.md`
- ✅ 包含所有必需的证据和链接

## 🚀 部署步骤

### 步骤 1：更新数据库配置（添加隐私政策和服务条款）

**方法 A - 通过管理后台（推荐）：**
1. 登录 https://your-domain.com 管理后台
2. 进入"设置" → "系统设置"（可能需要添加此设置页面）
3. 找到"法律文档"配置
4. 粘贴 `docs/privacy_policy.md` 内容到"隐私政策"字段
5. 粘贴 `docs/terms_of_service.md` 内容到"服务条款"字段
6. 点击"保存"

**方法 B - 通过 SSH 直接更新数据库：**
```bash
# 1. 连接到服务器
ssh root@YOUR_SERVER_IP

# 2. 找到数据库容器或直接连接数据库
# 如果使用 Docker:
docker exec -it <database-container> mysql -uroot -p<password> new-api

# 或如果直接安装:
mysql -uroot -p<password> new-api

# 3. 执行 SQL 更新（从本地复制 /tmp/update_legal_settings.sql 内容）
# 或手动执行 INSERT/UPDATE 语句
```

**方法 C - 使用提供的 SQL 文件：**
```bash
# 本地有 SQL 文件：/tmp/update_legal_settings.sql
# 需要手动上传到服务器并执行
```

### 步骤 2：部署前端更新

```bash
# 从本地执行（当前目录：yi-api）
cd web

# 构建前端（已完成）
bun run build

# 部署到服务器
# 方法 1：使用 rsync（需要先在服务器创建目录）
ssh root@YOUR_SERVER_IP "mkdir -p /your/deploy/path/web"
rsync -avz --delete dist/ root@YOUR_SERVER_IP:/your/deploy/path/web/dist/

# 方法 2：通过 Git 推送并在服务器拉取
cd ..
git add -A
git commit -m "feat: 添加隐私政策、服务条款和客户支持邮箱以满足 Creem 合规要求"
git push
# 然后在服务器上：
ssh root@YOUR_SERVER_IP
cd /your/deploy/path
git pull
cd web && bun run build
```

### 步骤 3：重启服务（如需）

```bash
# 在服务器上
ssh root@YOUR_SERVER_IP

# 如果使用 systemd
systemctl restart new-api

# 如果使用 docker-compose
cd /your/deploy/path && docker-compose restart new-api

# 如果直接运行
pkill -f new-api
cd /your/deploy/path && nohup ./new-api > /dev/null 2>&1 &
```

### 步骤 4：验证部署

1. **验证前端更新：**
   - 访问 https://your-domain.com
   - 滚动到页脚，确认显示以下链接：
     - ✅ 隐私政策
     - ✅ 服务条款
     - ✅ 客户支持（邮箱）

2. **验证法律文档可访问：**
   - 访问 https://your-domain.com/privacy
   - 访问 https://your-domain.com/user-agreement
   - 确认内容正确显示

3. **验证客户支持邮箱：**
   - 点击页脚的"客户支持"链接
   - 确认打开邮件客户端，收件人为：your-email@example.com

### 步骤 5：回复 Creem 团队

1. 获取 Store ID（从 Creem 管理后台）
2. 打开 `docs/CREEM_REPLY_EMAIL.md`
3. 替换 `[YOUR_STORE_ID]` 为实际的 Store ID
4. 复制邮件内容
5. 发送到：support@creem.io
6. （可选）附上页脚截图作为证据

## 📋 需要的截图

为 Creem 团队准备以下截图：

1. **网站页脚截图**
   - 显示"隐私政策"、"服务条款"、"客户支持"链接
   - URL: https://your-domain.com

2. **隐私政策页面截图**
   - URL: https://your-domain.com/privacy

3. **服务条款页面截图**
   - URL: https://your-domain.com/user-agreement

## 🔍 故障排查

**如果隐私政策/服务条款页面显示"加载失败"：**
1. 检查数据库配置是否已更新
2. 检查 API 端点 `/api/privacy-policy` 和 `/api/user-agreement` 是否返回内容
3. 重启服务确保配置生效

**如果页脚链接不显示：**
1. 确认前端已重新构建并部署
2. 清除浏览器缓存
3. 检查 `Footer.jsx` 文件是否正确更新

## 📝 相关文件

- `docs/privacy_policy.md` - 隐私政策原文
- `docs/terms_of_service.md` - 服务条款原文
- `docs/CREEM_REPLY_EMAIL.md` - 给 Creem 的回复邮件
- `docs/UPDATE_LEGAL_SETTINGS.md` - 详细的配置更新指南
- `web/src/components/layout/Footer.jsx` - 更新的页脚组件
- `/tmp/update_legal_settings.sql` - SQL 更新脚本

## ⚠️ 重要提醒

1. **必须先部署前端和数据库配置**，然后再回复 Creem 团队
2. **验证所有链接都能正常访问**后再发送邮件
3. **保存 Store ID** 以便在邮件中使用
4. **准备截图**作为证据附件

## 🎯 当前状态

- [x] 隐私政策文档已创建
- [x] 服务条款文档已创建
- [x] 前端 Footer 已更新
- [x] 前端已构建
- [ ] 数据库配置已更新（需手动执行）
- [ ] 前端已部署到生产
- [ ] 服务已重启
- [ ] 验证所有链接可访问
- [ ] 回复邮件已发送给 Creem
