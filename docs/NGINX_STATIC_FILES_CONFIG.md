# Nginx 静态文件服务配置 - 微信二维码图片访问

## 问题描述

yi-api 只服务编译时嵌入的静态文件（`web/dist`），无法访问运行时文件系统上的 `/www/wwwroot/api.ai80.vip/images/` 目录。

## 解决方案：配置 Nginx 单独服务 `/images` 目录

### 1. 编辑 Nginx 配置文件

找到你的站点配置文件，通常位于：
- `/etc/nginx/sites-available/api.ai80.vip.conf`
- `/etc/nginx/conf.d/api.ai80.vip.conf`
- `/www/server/panel/vhost/nginx/api.ai80.vip.conf` （宝塔面板）

### 2. 添加静态文件 location 配置

在 `server` 块中添加以下配置（**在反向代理配置之前**）：

```nginx
server {
    listen 80;
    server_name api.ai80.vip;

    # 1. 静态文件直接服务（优先级最高，放在最前面）
    location /images/ {
        alias /www/wwwroot/api.ai80.vip/images/;

        # 允许跨域访问（如果需要）
        add_header 'Access-Control-Allow-Origin' '*';
        add_header 'Access-Control-Allow-Methods' 'GET, OPTIONS';

        # 缓存配置（可选）
        expires 7d;
        add_header Cache-Control "public, immutable";

        # 安全：禁止目录浏览
        autoindex off;

        # 日志（可选，用于调试）
        access_log /var/log/nginx/api.ai80.vip_images_access.log;
        error_log /var/log/nginx/api.ai80.vip_images_error.log;
    }

    # 2. 反向代理到 yi-api（处理其他所有请求）
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 流式响应支持
        proxy_buffering off;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
    }

    # SSL 配置（如果有证书）
    # listen 443 ssl;
    # ssl_certificate /path/to/cert.pem;
    # ssl_certificate_key /path/to/key.pem;
}
```

### 3. 测试配置并重载 Nginx

```bash
# 测试配置文件语法
sudo nginx -t

# 如果显示 "syntax is ok" 和 "test is successful"，则重载配置
sudo nginx -s reload
# 或
sudo systemctl reload nginx
```

### 4. 验证访问

```bash
# 测试图片是否可访问
curl -I http://api.ai80.vip/images/wechat_qrcode.png

# 预期输出：
# HTTP/1.1 200 OK
# Content-Type: image/png
# Content-Length: 28519
```

在浏览器访问：
```
http://api.ai80.vip/images/wechat_qrcode.png
```

---

## 方法二：宝塔面板配置（GUI 方式）

如果您使用宝塔面板：

### 步骤 1: 进入站点设置

1. 登录宝塔面板
2. 网站 → 找到 `api.ai80.vip` → 设置

### 步骤 2: 添加反向代理规则

1. 点击"反向代理"标签
2. 确保有反向代理到 `http://127.0.0.1:3000`

### 步骤 3: 配置静态文件目录

1. 点击"配置文件"标签
2. 在 `location /` 之前添加：

```nginx
location /images/ {
    alias /www/wwwroot/api.ai80.vip/images/;
    expires 7d;
    add_header Cache-Control "public, immutable";
}
```

3. 点击"保存"

---

## 方法三：Docker 挂载静态文件目录

修改 `docker-compose.yml`，将 images 目录挂载到容器内的 `web/dist/images`：

```yaml
services:
  new-api:
    image: calciumion/new-api:latest
    container_name: new-api
    restart: always
    ports:
      - "3000:3000"
    volumes:
      - ./data:/data
      - ./logs:/app/logs
      # 新增：挂载 images 目录
      - /www/wwwroot/api.ai80.vip/images:/app/web/dist/images
    environment:
      # ... 其他配置
```

**重启服务**:
```bash
docker-compose down
docker-compose up -d
```

**注意**: 这种方法需要确保容器内的路径 `/app/web/dist/images` 存在且有正确的权限。

---

## 推荐配置（完整 Nginx 示例）

### HTTP + HTTPS 完整配置

```nginx
# HTTP 重定向到 HTTPS（生产环境推荐）
server {
    listen 80;
    server_name api.ai80.vip;
    return 301 https://$server_name$request_uri;
}

# HTTPS 配置
server {
    listen 443 ssl http2;
    server_name api.ai80.vip;

    # SSL 证书配置
    ssl_certificate /etc/nginx/ssl/api.ai80.vip/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/api.ai80.vip/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # 安全头
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;

    # 客户端请求体大小限制（用于文件上传）
    client_max_body_size 100M;

    # 静态文件服务（优先级最高）
    location /images/ {
        alias /www/wwwroot/api.ai80.vip/images/;

        # 跨域配置
        add_header 'Access-Control-Allow-Origin' '*' always;
        add_header 'Access-Control-Allow-Methods' 'GET, OPTIONS' always;
        add_header 'Access-Control-Allow-Headers' 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range' always;

        # 缓存
        expires 7d;
        add_header Cache-Control "public, immutable";

        # 安全
        autoindex off;

        # 支持 OPTIONS 预检请求
        if ($request_method = 'OPTIONS') {
            add_header 'Access-Control-Max-Age' 1728000;
            add_header 'Content-Type' 'text/plain; charset=utf-8';
            add_header 'Content-Length' 0;
            return 204;
        }
    }

    # API 反向代理
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;

        # 请求头
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        # 超时配置
        proxy_connect_timeout 60s;
        proxy_send_timeout 300s;
        proxy_read_timeout 300s;

        # 流式响应（SSE 支持）
        proxy_buffering off;
        proxy_cache off;
    }

    # WebSocket 支持（如果需要）
    location /ws/ {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_read_timeout 86400;
    }

    # 日志
    access_log /var/log/nginx/api.ai80.vip_access.log;
    error_log /var/log/nginx/api.ai80.vip_error.log;
}
```

---

## 文件权限检查

确保 Nginx 进程有权限读取 images 目录：

```bash
# 查看 Nginx 运行用户
ps aux | grep nginx

# 设置目录权限（假设 Nginx 用户为 www-data 或 nginx）
sudo chown -R www-data:www-data /www/wwwroot/api.ai80.vip/images/
sudo chmod -R 755 /www/wwwroot/api.ai80.vip/images/

# 或者（宝塔面板默认）
sudo chown -R www:www /www/wwwroot/api.ai80.vip/images/
sudo chmod -R 755 /www/wwwroot/api.ai80.vip/images/
```

---

## 故障排查

### 1. 检查 Nginx 配置语法

```bash
sudo nginx -t
```

### 2. 查看 Nginx 错误日志

```bash
sudo tail -f /var/log/nginx/error.log
```

### 3. 测试文件访问

```bash
# 直接访问文件（绕过 Nginx）
cat /www/wwwroot/api.ai80.vip/images/wechat_qrcode.png

# 检查文件权限
ls -la /www/wwwroot/api.ai80.vip/images/
```

### 4. 测试 Nginx 配置

```bash
# HTTP 请求测试
curl -I http://api.ai80.vip/images/wechat_qrcode.png

# HTTPS 请求测试
curl -I https://api.ai80.vip/images/wechat_qrcode.png

# 详细输出
curl -v http://api.ai80.vip/images/wechat_qrcode.png
```

### 5. 检查防火墙

```bash
# 检查端口是否开放
sudo netstat -tlnp | grep :80
sudo netstat -tlnp | grep :443

# 检查防火墙规则（CentOS/RHEL）
sudo firewall-cmd --list-all

# 检查防火墙规则（Ubuntu/Debian）
sudo ufw status
```

---

## 常见问题

### Q1: 配置后仍然 404？

**检查 location 顺序**: 确保 `location /images/` 在 `location /` **之前**。

Nginx 按顺序匹配，如果 `location /` 在前，会拦截所有请求。

### Q2: 403 Forbidden 错误？

**检查权限**:
```bash
sudo chmod 755 /www/wwwroot/api.ai80.vip/images/
sudo chmod 644 /www/wwwroot/api.ai80.vip/images/wechat_qrcode.png
```

**检查 SELinux**（CentOS/RHEL）:
```bash
sudo setenforce 0  # 临时禁用测试
```

### Q3: 如何支持其他静态文件目录？

添加更多 location 块：

```nginx
location /uploads/ {
    alias /www/wwwroot/api.ai80.vip/uploads/;
    expires 7d;
}

location /assets/ {
    alias /www/wwwroot/api.ai80.vip/assets/;
    expires 30d;
}
```

---

## 生产环境建议

1. **使用 HTTPS** - 保护用户隐私和数据安全
2. **启用 HTTP/2** - 提升性能
3. **配置缓存** - 减少服务器负载
4. **限制文件类型** - 防止恶意文件上传
5. **监控日志** - 及时发现异常访问

### 限制文件类型示例

```nginx
location /images/ {
    alias /www/wwwroot/api.ai80.vip/images/;

    # 只允许图片文件
    location ~* \.(jpg|jpeg|png|gif|webp|svg)$ {
        expires 7d;
        add_header Cache-Control "public, immutable";
    }

    # 拒绝其他文件类型
    location ~ {
        deny all;
    }
}
```

---

## 总结

**推荐方案**: 配置 Nginx 单独服务 `/images/` 目录

**关键配置**:
```nginx
location /images/ {
    alias /www/wwwroot/api.ai80.vip/images/;
    expires 7d;
    add_header Cache-Control "public, immutable";
}
```

**验证步骤**:
1. `sudo nginx -t` - 测试配置
2. `sudo nginx -s reload` - 重载 Nginx
3. `curl -I http://api.ai80.vip/images/wechat_qrcode.png` - 测试访问

配置完成后，图片即可通过 `http://api.ai80.vip/images/wechat_qrcode.png` 正常访问。
