# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Yi-API (New API) 是一个企业级多渠道AI API网关和中继服务平台。它统一管理并转发用户请求到30+家AI服务商（OpenAI、Claude、Gemini、DeepSeek、阿里云、百度等），提供配额管理、成本控制、渠道负载均衡等企业级功能。

**技术栈**:
- 后端: Go 1.25+ + Gin Framework
- 前端: React 18 + Vite + Semi Design (抖音)
- 数据库: MySQL/PostgreSQL/SQLite (GORM)
- 缓存: Redis (可选) + 内存缓存

## 常用命令

### 开发环境

**后端开发**:
```bash
# 安装依赖
go mod download

# 运行后端服务
go run main.go

# 调试模式 (需要设置环境变量)
GIN_MODE=debug go run main.go

# 查看环境配置
cp .env.example .env  # 首次运行需要创建环境配置
```

**前端开发**:
```bash
cd web

# 安装依赖 (使用 bun 或 npm)
bun install
# 或
npm install

# 启动开发服务器
bun run dev
# 或
npm run dev

# 构建生产版本
bun run build
# 或
npm run build

# 代码格式化
bun run lint:fix        # 自动修复格式问题
bun run eslint:fix      # 自动修复 ESLint 问题
```

**完整构建**:
```bash
# 使用 Makefile 构建整个项目 (前端 + 后端)
make all
```

### Docker 部署

详细的 Docker 部署说明请参考 [README.md](./README.md) 或 [官方文档](https://docs.newapi.pro/installation)。

**快速命令**:
```bash
# Docker Compose (推荐)
docker-compose up -d

# 查看日志
docker-compose logs -f new-api
```

### 数据库相关

**迁移和初始化**:
```bash
# 数据库会在首次启动时自动初始化和迁移
# GORM 会自动创建表结构，无需手动执行迁移脚本
```

**连接数据库** (用于调试):
```bash
# PostgreSQL
docker exec -it postgres psql -U root -d new-api

# MySQL (如果使用 MySQL)
docker exec -it mysql mysql -uroot -p123456 new-api
```

## 核心架构

### 分层架构

```
用户请求
    ↓
[Router层] router/         - 路由分配与中间件链
    ↓
[Middleware层] middleware/ - 认证、限流、缓存、日志
    ↓
[Controller层] controller/ - HTTP请求处理与业务流程组织
    ↓
[Service层] service/       - 核心业务逻辑实现
    ↓
[Relay层] relay/           - 多渠道适配器与API转发
    ↓
[Model层] model/           - 数据持久化与缓存管理
    ↓
[外部API]                  - 上游AI服务商
```

### 核心模块说明

**router/** - 路由层
- `main.go` - 路由入口与整合
- `api-router.go` - 后台管理API路由 (/api/*)
- `relay-router.go` - **API中继路由** (/v1/*, 对外提供OpenAI兼容API)
- `dashboard.go` - 前端页面路由
- `web-router.go` - 静态资源服务

**middleware/** - 中间件层 (17个文件)
- `auth.go` - Token/用户认证、权限验证 (RootAuth/AdminAuth/UserAuth)
- `rate-limit.go` - 全局/用户级别流量限制 (基于Redis)
- `model-rate-limit.go` - 模型级别限流
- `distributor.go` - **渠道分配与负载均衡** (核心组件)
- `cache.go` - 响应缓存
- `logger.go` - 请求日志记录

**relay/** - 中继转发层 (最复杂的模块)
- `channel/` - 30+个渠道适配器目录 (openai/, claude/, gemini/, deepseek/, ali/, baidu/等)
- `common/` - 中继公共逻辑
- `helper/` - 辅助工具 (模型映射、价格计算、请求验证)
- `*_handler.go` - 特定类型请求处理器:
  - `compatible_handler.go` - 通用文本请求 (OpenAI兼容)
  - `image_handler.go` - 图像生成
  - `embedding_handler.go` - 文本嵌入
  - `audio_handler.go` - 音频处理
  - `websocket.go` - WebSocket实时流

**Adaptor接口** - 核心抽象 (所有渠道适配器都实现此接口):
```go
type Adaptor interface {
    Init(info *RelayInfo)
    GetRequestURL() string
    SetupRequestHeader()
    ConvertOpenAIRequest() (any, error)      // 请求转换
    DoRequest() (any, error)                 // 执行请求
    DoResponse() (usage any, err)            // 响应处理
    GetModelList() []string
    GetChannelName() string
}
```

**controller/** - 控制器层 (46个文件)
按功能分类:
- 用户管理: user.go, passkey.go, twofa.go
- 渠道管理: channel.go, channel-test.go, channel-billing.go
- 模型管理: model.go, model_meta.go, model_sync.go
- 配额管理: token.go, billing.go, redemption.go
- 定价配置: pricing.go, ratio_config.go, ratio_sync.go
- 任务管理: task.go, task_video.go, midjourney.go
- 支付集成: topup.go, topup_stripe.go, topup_creem.go

**service/** - 服务层
- `channel_select.go` - **渠道智能选择** (自动分组、负载均衡)
- `quota.go` - **配额计算与扣费** (预扣费、补扣费、返还逻辑)
- `token_counter.go` - Token计数 (支持多种编码算法)
- `error.go` - 错误处理与重试
- `http.go` - HTTP客户端管理
- `webhook.go` - 事件通知

**model/** - 数据模型层 (GORM)
核心数据表:
- `User` - 用户账户 (quota配额管理、role权限、group分组)
- `Channel` - AI服务渠道 (type类型、key密钥、models支持的模型、weight权重)
- `Token` - API Token (userId、quota额度、modelLimits模型限制)
- `Ability` - 渠道-模型关系表 (channelId + modelName多对多)
- `Log` - 消耗日志表 (大表，记录所有API调用)
- `Pricing` - 模型定价表
- `Task` - 异步任务 (Suno音乐、Kling视频等)

**common/** - 公共模块 (35个文件)
- `crypto.go` - 加密解密
- `redis.go` - Redis客户端
- `database.go` - 数据库连接
- `limiter/` - 限流器实现 (Lua脚本)

**setting/** - 配置管理层
- `operation_setting/` - 运营配置 (通用设置、配额设置、支付配置)
- `ratio_setting/` - 定价倍率 (模型倍率、分组倍率、缓存倍率)
- `model_setting/` - 模型专用设置

**constant/** - 常量定义
- APIType枚举 (OpenAI/Claude/Gemini等)
- ChannelType枚举
- TaskPlatform枚举

### 关键工作流

**API请求处理完整流程**:

```
POST /v1/chat/completions
    ↓
[TokenAuth] 认证
  → 验证 API Token
  → 获取用户配额信息
    ↓
[RateLimit] 限流检查
  → Redis流控验证
    ↓
[RequestValidation] 请求验证
  → 检查模型是否存在
  → 验证请求参数
    ↓
[ModelMapping] 模型映射
  → channel.model_mapping 转换模型名称
    ↓
[ChannelSelect] 渠道选择 (service/channel_select.go)
  → 按group查找可用渠道
  → 按model_name过滤
  → 按weight加权随机选择
  → 返回最终Channel
    ↓
[PreConsumeQuota] 预扣费 (service/pre_consume_quota.go)
  → 估算token数量
  → 计算预计配额
  → 冻结用户/Token余额
    ↓
[AdaptorSelect] 选择适配器
  → GetAdaptor(channel.type) 获取对应渠道适配器
    ↓
[RequestConvert] 请求转换 (Adaptor.ConvertOpenAIRequest)
  → OpenAI格式 → 目标渠道格式 (Claude/Gemini等)
  → 应用参数覆盖 (ParamOverride)
  → 应用请求头覆盖 (HeaderOverride)
    ↓
[DoRequest] 执行转发 (Adaptor.DoRequest)
  → 设置请求头 (API Key、Auth等)
  → 处理流式响应 (SSE)
  → 错误处理与渠道故障标记
    ↓
[DoResponse] 响应处理 (Adaptor.DoResponse)
  → 流式数据转发
  → 实时解析Usage信息
  → 返回最终Usage
    ↓
[PostConsumeQuota] 补扣费 (service/quota.go)
  → 计算实际消耗
  → quotaDelta = actual - pre
  → 调整配额 (返还或继续扣费)
    ↓
[RecordLog] 记录日志
  → 用户消耗日志
  → 渠道使用统计
    ↓
返回响应给客户端
```

### 核心设计模式

**1. 渠道选择与负载均衡**

```go
// service/channel_select.go
// 策略:
// - 权重模式: 按weight字段加权随机选择
// - 自动分组模式: 根据用户配置自动选择合适分组
// - 故障转移: 失败时自动重试其他渠道

func CacheGetRandomSatisfiedChannel(
    group string,    // 渠道分组
    model string,    // 模型名称
    retry int        // 重试次数
) (*Channel, error)
```

**2. 多Key管理机制**

Channel支持多个API Key轮询使用:
```
Channel.ChannelInfo {
    IsMultiKey: true              # 启用多Key功能
    MultiKeySize: 5               # 配置5个Key
    MultiKeyMode: "polling"       # 轮询模式 (或"random"随机)
    MultiKeyStatusList: {
        0: enabled,
        1: disabled,              # 某个Key失败后自动禁用
        2: enabled,
        ...
    }
    MultiKeyPollingIndex: 3       # 当前轮询索引
}
```

**3. 配额与成本控制**

**预扣费机制** (Pre-consume Quota):
- 请求前先冻结预估配额，防止超额使用
- 请求完成后根据实际消耗补扣费或返还多扣部分

**配额计算公式**:
```
quota = (
    promptTokens × 1.0
    + cachedTokens × cacheRatio          # 缓存内容折扣
    + completionTokens × completionRatio  # 输出token倍率
    + imageTokens × imageRatio
) × modelRatio × groupRatio × quotaPerUnit

# 或使用固定价格模式:
quota = modelPrice × quotaPerUnit × groupRatio
```

**倍率系统**:
- **模型倍率**: 不同模型的相对成本
- **分组倍率**: 不同用户组的差异化定价
- **缓存倍率**: Cache命中时的优惠
- **补全倍率**: 输出相对输入的成本倍数

**4. 流式处理 (SSE)**

```
请求 → 流式响应
├─ 设置 Content-Type: text/event-stream
├─ 逐块读取、转换、转发上游数据
├─ 实时解析Usage (计算实际消耗)
└─ 连接关闭时补扣费
```

## 重要配置

### 环境变量

关键环境变量 (配置在 `.env` 文件或 docker-compose.yml):

```bash
# 数据库配置
SQL_DSN=postgresql://root:123456@localhost:5432/new-api
# 或 MySQL
SQL_DSN=root:123456@tcp(localhost:3306)/new-api

# Redis配置 (可选，启用缓存)
REDIS_CONN_STRING=redis://localhost:6379

# 服务配置
PORT=3000
GIN_MODE=release                    # release/debug
TZ=Asia/Shanghai                    # 时区

# 集群配置
SESSION_SECRET=random_string        # 多机部署时必须设置相同值
CRYPTO_SECRET=random_string         # Redis模式下必须设置
IsMasterNode=true                   # 是否为主节点

# 缓存与同步
MEMORY_CACHE_ENABLED=true           # 启用内存缓存
SYNC_FREQUENCY=60                   # 缓存同步频率(秒)

# 性能配置
STREAMING_TIMEOUT=300               # 流式超时时间(秒)
BATCH_UPDATE_ENABLED=true           # 批量更新开关

# 日志配置
ERROR_LOG_ENABLED=true              # 错误日志开关

# 调试配置
ENABLE_PPROF=true                   # 启用性能分析(仅开发环境)
```

### 多机部署注意事项

**必须配置**:
- `SESSION_SECRET` - 保证登录状态一致
- `CRYPTO_SECRET` - 使用Redis时保证数据加密一致
- `REDIS_CONN_STRING` - 共享缓存

**主从节点**:
- 主节点 (`IsMasterNode=true`): 处理缓存同步、异步任务更新
- 从节点: 只处理API请求

## 开发指南

### 添加新的AI渠道适配器

1. **在 `relay/channel/` 创建新目录**:
```bash
mkdir relay/channel/newprovider
```

2. **实现 Adaptor 接口** (`relay/channel/newprovider/adaptor.go`):
```go
package newprovider

import "github.com/QuantumNous/new-api/relay/common"

type Adaptor struct {
    common.BaseAdaptor
    // 添加渠道特定字段
}

func (a *Adaptor) Init(info *common.RelayInfo) {
    a.BaseAdaptor.Init(info)
}

func (a *Adaptor) GetRequestURL() string {
    // 返回API端点URL
    return "https://api.newprovider.com/v1/chat"
}

func (a *Adaptor) SetupRequestHeader() error {
    // 设置认证头等
    a.Request.Header.Set("Authorization", "Bearer " + a.ChannelKey)
    return nil
}

func (a *Adaptor) ConvertOpenAIRequest() (any, error) {
    // 将OpenAI格式请求转换为目标格式
    request := a.GetOpenAIRequest()
    // ... 转换逻辑
    return convertedRequest, nil
}

func (a *Adaptor) DoRequest() (any, error) {
    // 执行HTTP请求
    return a.BaseAdaptor.DoRequestHelper(a)
}

func (a *Adaptor) DoResponse() (*common.Usage, error) {
    // 处理响应并返回Usage信息
    // ... 响应处理逻辑
    return &common.Usage{
        PromptTokens:     promptTokens,
        CompletionTokens: completionTokens,
        TotalTokens:      totalTokens,
    }, nil
}

func (a *Adaptor) GetModelList() []string {
    return []string{"model-1", "model-2"}
}

func (a *Adaptor) GetChannelName() string {
    return "NewProvider"
}
```

3. **注册适配器** (在 `relay/common/adaptor.go` 或相应的工厂函数):
```go
case constant.ChannelTypeNewProvider:
    return &newprovider.Adaptor{}, nil
```

4. **添加渠道类型常量** (`constant/channel.go`):
```go
const (
    // ...
    ChannelTypeNewProvider = 99
)

var ChannelType2Name = map[int]string{
    // ...
    ChannelTypeNewProvider: "newprovider",
}
```

### 修改配额计算逻辑

配额计算主要在 `service/quota.go` 和 `service/pre_consume_quota.go`:

```go
// service/quota.go

// 修改配额计算公式
func CalculateQuota(
    promptTokens int,
    completionTokens int,
    cachedTokens int,
    modelRatio float64,
    groupRatio float64,
    cacheRatio float64,
    completionRatio float64,
) int64 {
    quota := float64(promptTokens)
    quota += float64(cachedTokens) * cacheRatio
    quota += float64(completionTokens) * completionRatio
    quota *= modelRatio * groupRatio
    return int64(quota)
}
```

### 添加新的中间件

在 `middleware/` 目录创建新的中间件文件:

```go
package middleware

import "github.com/gin-gonic/gin"

func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理

        c.Next()  // 执行后续处理器

        // 后置处理
    }
}
```

然后在 `router/` 中注册:

```go
router.Use(middleware.CustomMiddleware())
```

### 代码风格与规范

**遵循 Go 标准规范**:
- 使用 `gofmt` 格式化代码
- 遵循 [Effective Go](https://go.dev/doc/effective_go) 指南
- 使用驼峰命名法 (CamelCase)
- 导出函数/结构体首字母大写，私有的首字母小写

**错误处理**:
```go
// 好的做法
result, err := SomeFunction()
if err != nil {
    return nil, fmt.Errorf("操作失败: %w", err)
}

// 避免忽略错误
result, _ := SomeFunction()  // 不推荐
```

**日志记录**:
```go
import "github.com/QuantumNous/new-api/common"

// 使用项目的日志函数
common.SysLog("系统日志信息")
common.SysError("错误信息")
common.FatalLog("严重错误")
```

## 故障排查

### 常见问题

**1. 数据库连接失败**
```bash
# 检查数据库配置
echo $SQL_DSN

# 测试数据库连接
# PostgreSQL
docker exec -it postgres psql -U root -d new-api -c "SELECT 1"
# MySQL
docker exec -it mysql mysql -uroot -p123456 -e "SELECT 1"
```

**2. Redis连接失败**
```bash
# 检查Redis配置
echo $REDIS_CONN_STRING

# 测试Redis连接
docker exec -it redis redis-cli ping
```

**3. 渠道请求失败**
- 检查 Channel 配置中的 API Key 是否正确
- 检查 Channel 状态是否为启用
- 查看日志中的详细错误信息 (日志位于 `./logs/` 或容器内 `/app/logs/`)
- 使用"渠道测试"功能测试单个渠道

**4. 前端编译失败**
```bash
# 清理缓存重新安装
cd web
rm -rf node_modules
bun install  # 或 npm install
bun run build
```

**5. Token消耗异常**
- 检查模型倍率配置 (设置 → 运营设置 → 倍率设置)
- 检查用户分组倍率
- 查看消耗日志中的详细计算过程

### 日志查看

**Docker Compose 部署**:
```bash
# 查看实时日志
docker-compose logs -f new-api

# 查看最近100行
docker-compose logs --tail 100 new-api

# 查看错误日志
docker exec -it new-api cat /app/logs/error.log
```

**本地开发**:
```bash
# 日志会输出到终端
go run main.go

# 或查看日志文件
tail -f ./logs/error.log
```

## 测试

**注意**: 项目目前没有后端单元测试。测试主要通过以下方式进行：

### 手动测试

**使用内置Playground**:
1. 登录管理后台
2. 进入 "Playground" 页面
3. 选择模型和渠道进行测试

**使用 curl 测试API**:
```bash
# 聊天接口测试
curl http://localhost:3000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-xxxxx" \
  -d '{
    "model": "gpt-4",
    "messages": [{"role": "user", "content": "Hello"}]
  }'

# 流式响应测试
curl http://localhost:3000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-xxxxx" \
  -d '{
    "model": "gpt-4",
    "messages": [{"role": "user", "content": "Hello"}],
    "stream": true
  }'
```

### 渠道测试

在管理后台使用"渠道测试"功能:
1. 进入 "渠道" 页面
2. 点击某个渠道的"测试"按钮
3. 查看测试结果和响应时间

## 性能优化

- 启用 Redis 缓存: `REDIS_CONN_STRING=redis://localhost:6379`
- 启用内存缓存: `MEMORY_CACHE_ENABLED=true`
- 启用批量更新: `BATCH_UPDATE_ENABLED=true`
- 生产环境建议使用 PostgreSQL 或 MySQL 而非 SQLite
- 可分离日志数据库: `LOG_SQL_DSN=...`

## 相关资源

- **官方文档**: https://docs.newapi.pro/
- **GitHub仓库**: https://github.com/Calcium-Ion/new-api
- **问题反馈**: https://github.com/Calcium-Ion/new-api/issues
- **上游项目**: [One API](https://github.com/songquanpeng/one-api)