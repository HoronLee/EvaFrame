# EvaFrame - Go Web 开发框架

一个现代化的 Go Web 开发框架，集成了 Gin、GORM、Wire 依赖注入、JWT 认证等核心组件。

## 技术栈

- **Web框架**: Gin - 高性能HTTP Web框架
- **ORM**: GORM - Go语言ORM库
- **依赖注入**: Wire - 编译时依赖注入
- **配置管理**: Viper - 多格式配置文件支持
- **日志**: Zap - 高性能结构化日志
- **JWT**: golang-jwt - JWT令牌认证
- **验证**: validator - 请求数据验证
- **命令行**: Cobra - 强大的CLI应用框架

## 项目结构

```bash
evaframe/
├── cmd/                   # 命令行接口
│   ├── root.go            # 根命令
│   ├── serve.go           # 服务器启动命令
│   └── migrate.go         # 数据库迁移命令
├── config/                # 配置文件
│   └── config.yaml        # 主配置文件
├── internal/              # 内部代码
│   ├── app/               # 应用程序入口
│   ├── dao/               # 数据访问层
│   ├── handler/           # HTTP处理器
│   ├── models/            # 数据模型
│   └── service/           # 业务逻辑层
└── pkg/                   # 公共包
    ├── config/            # 配置管理
    ├── database/          # 数据库连接
    ├── jwt/               # JWT认证
    ├── logger/            # 日志管理
    ├── middleware/        # 中间件
    ├── response/          # 响应处理
    └── validator/         # 数据验证
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置数据库

编辑 `config/config.yaml` 文件，配置您的数据库连接：

```yaml
database:
  dsn: "username:password@tcp(localhost:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local"
```

### 3. 生成依赖注入代码

```bash
make generate
# 或者
wire gen ./internal/app
```

### 4. 运行数据库迁移

```bash
make migrate
```

### 5. 运行应用

```bash
make run
# 或者开发环境一键启动（迁移+运行）
make dev
```

## 数据库迁移

EvaFrame 内置了 GORM 自动迁移功能，可以自动创建和更新数据库表结构：

### 运行迁移
```bash
# 使用 Makefile
make migrate

# 使用自定义配置文件
go run main.go migrate --config /path/to/config.yaml
```

## 编码须知

### 编码顺序

基于 Kratos 风格的分层架构，实现新需求的标准编码顺序：

#### 1. 数据模型层
在 `internal/models/` 目录下定义领域实体：

```go
// internal/models/your_model.go
type YourModel struct {
  ID        uint      `gorm:"primarykey" json:"id"`
  Name      string    `json:"name"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

```
然后在 `cmd/migrate.go` 文件中添加到 AutoMigrate 列表：

```go
err = db.AutoMigrate(
    &models.User{},
    &models.YourModel{}, // 添加新模型
)
```

#### 2. Service 层定义 DAO 接口
在 Service 层明确需要哪些数据操作方法：

```go
// internal/service/your_model.go
type YourModelDAO interface {
  Create(model *models.YourModel) error
  GetByID(id uint) (*models.YourModel, error)
  List(offset, limit int) ([]*models.YourModel, int64, error)
  // 根据业务需求定义其他方法
}
```

#### 3. DAO 层实现接口
在 `internal/dao/gorm/` 目录下实现 Service 层定义的接口：

```go
// internal/dao/gorm/your_model.go
type YourModelDAOImpl struct {
  db *gorm.DB
}

func NewYourModelDAO(db *gorm.DB) service.YourModelDAO {
  return &YourModelDAOImpl{db: db}
}

// 实现接口方法...
```

#### 4. Service 层实现业务逻辑
编写纯业务逻辑方法，使用基本类型参数：

```go
// internal/service/your_model.go
func (s *YourModelService) CreateYourModel(name string) (*models.YourModel, error) {
    // 业务逻辑处理
    model := &models.YourModel{Name: name}
    return model, s.yourModelDAO.Create(model)
}
```

#### 5. Handler 层处理 HTTP
在 `internal/handler/` 目录下定义请求/响应结构体和处理器：

```go
// internal/handler/your_model.go
type CreateYourModelRequest struct {
    Name string `json:"name" validate:"required"`
}

func (h *YourModelHandler) Create(c *gin.Context) {
    var req CreateYourModelRequest
    // HTTP 协议处理、数据验证
    // 调用 Service 层业务逻辑
    // 返回响应
}
```

#### 6. 注册路由和依赖注入
- 在 Handler 中注册路由
- 在对应的 `gorm.go`、`service.go`、`handler.go` 文件中更新 Wire ProviderSet
- 在 `internal/app/app.go` 中注册新的 Handler

```go
type Application struct {
	YourModelHandler *handler.YourModelHandler
}

func NewApplication(
	yourModelHandler *handler.YourModelHandler,
) *Application {
	yourModelHandler.RegisterRoutes(router, authMiddleware)
	return &Application{
		YourModelHandler: userHandler,
	}
}
```

#### 7. 生成 Wire 代码和数据库迁移
```bash
# 生成依赖注入代码
make gen.wire

# 更新数据库迁移（在 cmd/migrate.go 中添加新模型）
make migrate
```

#### 8. 启动服务验证
```bash
# 开发环境一键启动
make dev
```

### 架构原则

- **Handler 层**：负责 HTTP 协议处理、请求验证、响应格式化
- **Service 层**：纯业务逻辑，定义 DAO 接口，不依赖外部框架
- **DAO 层**：实现 Service 层定义的接口，专注数据访问
- **依赖倒置**：Service 层定义接口，DAO 层实现接口
- **单一职责**：每层只关注自己的核心职责

这种架构确保了代码的可测试性、可维护性和可扩展性。
## CLI 命令

EvaFrame 提供了以下命令行工具：

- `serve` - 启动 Web 服务器
- `migrate` - 运行数据库迁移
- `--config` - 指定配置文件路径（全局选项）

```bash
# 查看所有可用命令
evaframe --help

# 查看特定命令帮助
evaframe migrate --help
```

## API 接口

### 用户注册
```bash
POST /api/v1/register
Content-Type: application/json

{
  "name": "张三",
  "email": "zhangsan@example.com",
  "password": "123456"
}
```

### 用户登录
```bash
POST /api/v1/login
Content-Type: application/json

{
  "email": "zhangsan@example.com",
  "password": "123456"
}
```

### 获取用户信息（需要JWT认证）
```bash
GET /api/v1/profile
Authorization: Bearer <your-jwt-token>
```

### 获取用户列表（需要JWT认证）
```bash
GET /api/v1/users?offset=0&limit=10
Authorization: Bearer <your-jwt-token>
```

## 可用命令

使用 Makefile 命令：

- `make build` - 编译应用程序
- `make run` - 运行服务器
- `make migrate` - 运行数据库迁移
- `make dev` - 开发环境快速启动（迁移+运行）
- `make clean` - 清理编译产物
- `make tidy` - 整理Go模块
- `make fmt` - 格式化代码
- `make generate` - 生成Wire代码
- `make help` - 显示帮助信息

## 配置说明

主配置文件位于 `config/config.yaml`：

```yaml
server:
  port: 8080              # 服务器端口
  mode: "debug"           # 运行模式: debug/release

database:
  type: "mysql" # 可选值: "mysql" 或 "sqlite"
  dsn: "..."              # 数据库连接字符串

jwt:
  secret: "..."           # JWT密钥

logger:
  level: "debug"          # 日志级别
  log_path: "./logs/app.log"  # 日志文件路径

dev_choice:
  dao: "gorm"             # DAO实现选择: gorm/gormgen
```

## 开发特性

- **依赖注入**: 使用 Wire 实现编译时依赖注入，确保类型安全
- **数据库迁移**: 内置 GORM 自动迁移，支持表结构自动创建和更新
- **配置热更新**: 支持配置文件热更新，无需重启服务
- **结构化日志**: 使用 Zap 提供高性能结构化日志
- **JWT认证**: 内置JWT中间件，支持用户认证
- **数据验证**: 使用 validator 进行请求数据验证
- **优雅关闭**: 支持服务器优雅关闭
- **命令行工具**: 基于 Cobra 的强大命令行接口

## 许可证

MIT License