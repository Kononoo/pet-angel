GOHOSTOS:=$(shell go env GOHOSTOS)                    # 本机 OS 名称（如 darwin、linux、windows）
GOPATH:=$(shell go env GOPATH)                        # Go 的 GOPATH
VERSION=$(shell git describe --tags --always)         # 使用 git tag/commit 作为版本

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto") # internal 下全部 proto
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")           # api 下全部 proto
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto) # 非 Windows 直接使用系统 find
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: init
# init: 一次性安装/更新本项目代码生成所需的工具链，包括：
# - protoc 的 Go 插件（protoc-gen-go、protoc-gen-go-grpc）
# - Kratos CLI 与 HTTP 协议生成器（protoc-gen-go-http）
# - OpenAPI 生成器（protoc-gen-openapi）
# - 依赖注入生成器（wire）
# 适用于新环境初始化或工具升级；执行后确保 $(GOPATH)/bin 已加入 PATH。
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: config
# config: 编译 `internal/` 目录下的 proto（如 `internal/conf/conf.proto`），
# 生成对应的 Go 类型（*.pb.go）到 `internal`，通常用于应用配置对象的更新。
# 当你调整配置的 proto 结构后，运行本目标以同步生成代码，并相应更新 `configs/config.yaml`。
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
  	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# api: 编译 `api/` 目录下的服务与消息定义 proto，生成：
# - *.pb.go（消息/枚举）
# - *_grpc.pb.go（gRPC 接口/客户端桩）
# - *_http.pb.go（Kratos HTTP 路由映射）
# 同时导出 OpenAPI 规范（openapi.yaml）便于接口文档对接。
# 修改对外 API（rpc/消息）后，先运行本目标，再在 `internal/service` 等处实现新增方法。
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
  	       --go_out=paths=source_relative:./api \
  	       --go-http_out=paths=source_relative:./api \
  	       --go-grpc_out=paths=source_relative:./api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

.PHONY: build
# build: 构建当前项目到 `bin/` 目录，使用 `-ldflags` 将 Git 版本信息注入 `main.Version`。
# 用于本地产出可执行文件或在 CI 中打包工件。
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate: 统一执行代码生成与依赖整理：
# - `go generate ./...` 扫描并执行源码中的 //go:generate 指令（如果存在）
# - `go mod tidy` 整理 go.mod / go.sum
# 通常在 `make api`/`make config` 之后执行，确保依赖图干净。
generate:
	go generate ./...
	go mod tidy

.PHONY: all
# all: 一键执行常用的生成链路（api → config → generate），
# 适合在修改 proto 后快速同步所有生成产物与依赖。
all:
	make api;
	make config;
	make generate;

# help: 显示可用的 make 目标与其简要说明（从各目标上方的注释中提取并渲染），
# 便于探索和记忆命令。默认目标即为 help，不带任何参数直接 `make` 即可查看。
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

# 默认目标：当不指定任何 target 时，展示 help 列表与简述。
.DEFAULT_GOAL := help
