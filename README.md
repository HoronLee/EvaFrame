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

```
evaframe/
├── cmd/                    # 命令行接口
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
├── pkg/                   # 公共包
│   ├── config/            # 配置管理
│   ├── database/          # 数据库连接
│   ├── jwt/               # JWT认证
│   ├── logger/            # 日志管理
│   ├── middleware/        # 中间件
│   ├── response/          # 响应处理
│   └── validator/         # 数据验证
└── tools/                 # 工具
    └── gormgen/           # GORM-Gen代码生成（已弃用）
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
# 或者
./bin/evaframe migrate
```

### 5. 运行应用

```bash
make run
# 或者开发环境一键启动（迁移+运行）
make dev
```

服务器将在 http://localhost:8080 启动

## 数据库迁移

EvaFrame 内置了 GORM 自动迁移功能，可以自动创建和更新数据库表结构：

### 运行迁移
```bash
# 使用 Makefile
make migrate

# 直接使用二进制文件
./bin/evaframe migrate

# 使用自定义配置文件
./bin/evaframe migrate --config /path/to/config.yaml
```

### 添加新模型
在 `internal/models/` 目录下创建新的模型文件，然后在 `cmd/migrate.go` 文件中添加到 AutoMigrate 列表：

```go
err = db.AutoMigrate(
    &models.User{},
    &models.YourNewModel{}, // 添加新模型
)
```

## CLI 命令

EvaFrame 提供了以下命令行工具：

- `serve` - 启动 Web 服务器
- `migrate` - 运行数据库迁移
- `--config` - 指定配置文件路径（全局选项）

```bash
# 查看所有可用命令
./bin/evaframe --help

# 查看特定命令帮助
./bin/evaframe migrate --help
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