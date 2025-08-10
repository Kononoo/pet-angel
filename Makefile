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
# init: 安装/更新代码生成工具链
# - protoc 插件: protoc-gen-go、protoc-gen-go-grpc
# - Kratos CLI 与 HTTP 协议生成器
# - OpenAPI 生成器
# - Wire 依赖注入生成器
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: api
# api: 生成 API 代码（pb.go / *_grpc.pb.go / *_http.pb.go）
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:./api \
	       --go-grpc_out=paths=source_relative:./api \
	       --go-http_out=paths=source_relative:./api \
	       $(API_PROTO_FILES)

.PHONY: openapi
# openapi: 导出 OpenAPI 文档至 openapi.yaml 方便前端联调与文档浏览
openapi:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)
	@echo "OpenAPI generated: ./openapi.yaml"

.PHONY: service
# service: 基于 proto 生成服务层骨架（仅生成缺失文件，不覆盖已实现代码）
# 示例：kratos proto server api/community/v1/community.proto -t internal/service
service:
	kratos proto server api/auth/v1/auth.proto -t internal/service || true
	kratos proto server api/avatar/v1/avatar.proto -t internal/service || true
	kratos proto server api/community/v1/community.proto -t internal/service || true
	kratos proto server api/user/v1/user.proto -t internal/service || true
	kratos proto server api/message/v1/message.proto -t internal/service || true
	kratos proto server api/upload/v1/upload.proto -t internal/service || true

.PHONY: config
# config: internal 下的 proto（例如 conf）生成到 internal 目录
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
   	   --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: wire
# wire: 生成依赖注入代码（cmd/pet-angel/wire_gen.go）
wire:
	cd cmd/pet-angel && wire

.PHONY: generate
# generate: 统一执行 generate 与依赖整理
generate:
	go generate ./...
	go mod tidy

.PHONY: run
# run: 以本地配置启动服务
run:
	go run cmd/pet-angel/main.go cmd/pet-angel/wire_gen.go -conf ./configs

.PHONY: build
# build: 构建可执行文件到 bin/
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/pet-angel ./cmd/pet-angel

.PHONY: test
# test: 运行项目全部测试
test:
	go test -v ./...

.PHONY: clean
# clean: 清理生成产物
clean:
	rm -rf bin/
	rm -f cmd/pet-angel/wire_gen.go

.PHONY: all
# all: 常用完整链路（生成 API/服务骨架/OpenAPI，wire 注入，并整理依赖）
all:
	make api
	make openapi
	make wire
	make generate

.PHONY: help
# help: 显示可用的 make 目标与其简要说明
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

.DEFAULT_GOAL := help
