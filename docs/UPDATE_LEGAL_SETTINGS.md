# 更新隐私政策和服务条款指南

本文档说明如何更新 Yi-API 的隐私政策和服务条款配置。

## 方法一：通过管理后台（推荐）

### 步骤：

1. 登录管理后台 https://your-domain.com
2. 进入"设置" → "系统设置"
3. 找到"法律文档"部分
4. 更新以下两个字段：
   - **用户协议 (Terms of Service)**: 粘贴 `docs/terms_of_service.md` 的内容
   - **隐私政策 (Privacy Policy)**: 粘贴 `docs/privacy_policy.md` 的内容
5. 点击"保存"

## 方法二：通过数据库直接更新

### PostgreSQL/MySQL:

```bash
# 连接到数据库
mysql -u root -p new-api

# 或 PostgreSQL:
psql -U root -d new-api

# 执行更新 SQL (从 /tmp/update_legal_settings.sql)
source /tmp/update_legal_settings.sql;
```

### 使用 Docker:

```bash
# 复制 SQL 文件到容器
docker cp /tmp/update_legal_settings.sql new-api-db:/tmp/

# 执行 SQL
docker exec -it new-api-db mysql -uroot -p123456 new-api < /tmp/update_legal_settings.sql
```

## 方法三：通过 API 更新

使用 API 端点更新配置：

```bash
# 获取管理员 Token (从管理后台)
TOKEN="your-admin-token-here"

# 准备 JSON 数据
cat > /tmp/legal_update.json << 'EOF'
{
  "key": "legal",
  "value": {
    "user_agreement": "... (完整内容见 docs/terms_of_service.md) ...",
    "privacy_policy": "... (完整内容见 docs/privacy_policy.md) ..."
  }
}
EOF

# 发送更新请求
curl -X PUT https://your-domain.com/api/option \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d @/tmp/legal_update.json
```

## 验证更新

更新后，访问以下链接验证：

- 隐私政策: https://your-domain.com/privacy
- 服务条款: https://your-domain.com/user-agreement

## 注意事项

1. 更新后需要重启服务使配置生效（如果使用数据库直接更新）
2. 确保备份原有配置
3. 客户支持邮箱已设置为: your-email@example.com

## 完整文档内容

完整的隐私政策和服务条款内容保存在：
- `docs/privacy_policy.md`
- `docs/terms_of_service.md`
