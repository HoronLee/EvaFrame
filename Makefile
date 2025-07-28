# ==============================================================================
# 项目变量
# ==============================================================================
APP_NAME := evaframe
MAIN_FILE := ./main.go
BINARY_PATH := ./bin/
BINARY_NAME := $(APP_NAME)
TARGET := $(BINARY_PATH)$(BINARY_NAME)
CONFIG_FILE := ./config/config.yaml

# Go 相关变量
GO := go
GO_BUILD := $(GO) build
GO_RUN := $(GO) run
GO_CLEAN := $(GO) clean
GO_FMT := gofmt -s -w
GO_TIDY := $(GO) mod tidy
GO_GENERATE := $(GO) generate

# ==============================================================================
# 命令目标
# ==============================================================================
.PHONY: all build run clean tidy fmt generate gen.wire migrate help

# 默认命令
all: build

# 编译应用程序
build:
	@echo "正在编译应用程序..."
	@mkdir -p $(BINARY_PATH)
	@$(GO_BUILD) -o $(TARGET) $(MAIN_FILE)
	@echo "二进制文件已生成: $(TARGET)"

# 运行应用程序
run:
	@echo "正在运行 Web 服务器..."
	@$(GO_RUN) $(MAIN_FILE) serve

# 运行数据库迁移
migrate: build
	@echo "正在运行数据库迁移..."
	@$(TARGET) migrate --config $(CONFIG_FILE)

# 开发环境快速启动 (迁移 + 运行)
dev: migrate
	@echo "正在启动开发环境..."
	@$(TARGET) serve --config $(CONFIG_FILE)

# 清理编译产物
clean:
	@echo "正在清理编译产物..."
	@rm -rf $(BINARY_PATH)
	@$(GO_CLEAN) -cache -testcache

# 整理 go.mod 文件
tidy:
	@echo "正在整理 Go 模块..."
	@$(GO_TIDY)

# 格式化代码
fmt:
	@echo "正在格式化 Go 代码..."
	@$(GO_FMT) ./...

# 生成 Wire 依赖注入代码
gen.wire:
	@echo "正在生成 Wire 依赖注入代码..."
	@wire gen ./internal/app

# 一键生成所有代码
generate: gen.wire
	@echo "所有代码已生成完毕。"

# 帮助信息
help:
	@echo ""
	@echo "用法: make [target]"
	@echo ""
	@echo "可用目标:"
	@echo "  all         (默认) 编译应用程序"
	@echo "  build       编译应用程序的二进制文件"
	@echo "  run         运行 Web 服务器"
	@echo "  migrate     运行数据库迁移"
	@echo "  dev         开发环境快速启动 (迁移 + 运行)"
	@echo "  clean       移除编译产物"
	@echo "  tidy        整理 Go 模块依赖"
	@echo "  fmt         格式化项目中的 Go 源代码"
	@echo "  generate    生成 Wire 依赖注入代码"
	@echo "  gen.wire    仅生成 Wire 依赖注入代码"
	@echo ""