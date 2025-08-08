# Pet Angel

宠物天使 - 一个专注于宠物全生命周期的社交应用

## 项目简介

Pet Angel是一个创新的宠物社交应用，致力于为宠物主人提供全生命周期的陪伴服务。从宠物的日常记录到社区交流，从虚拟形象互动到情感陪伴，我们希望通过技术的力量，让每一个宠物主人都能获得更好的养宠体验。

### 核心特色

- **虚拟形象系统**: 基于AI技术创建个性化宠物虚拟形象
- **社区互动**: 宠物主人社区，分享养宠经验和快乐
- **情感陪伴**: 宠物去世后的情感支持和陪伴
- **道具互动**: 丰富的道具系统，增强与虚拟宠物的互动
- **小纸条功能**: AI生成的个性化小纸条，提供情感慰藉

## 项目架构

### 技术栈

- **后端框架**: Kratos (Go微服务框架)
- **数据库**: MySQL 8.0+ (主数据库)
- **缓存**: Redis 6.0+ (缓存层)
- **API协议**: gRPC + HTTP (双协议支持)
- **依赖注入**: Wire (Google Wire)
- **配置管理**: 多环境配置支持

### 分层架构

```
┌─────────────────┐
│   API Layer     │  ← gRPC/HTTP API
├─────────────────┤
│  Service Layer  │  ← 业务服务层
├─────────────────┤
│   Biz Layer     │  ← 业务逻辑层
├─────────────────┤
│   Data Layer    │  ← 数据访问层
├─────────────────┤
│  Infrastructure │  ← MySQL/Redis/File Storage
└─────────────────┘
```

### 核心功能模块

#### 1. 用户系统 (User)
- **功能**: 用户注册、登录、信息管理、关注系统
- **API**: `/v1/user/*`
- **核心表**: `users`, `pets`, `user_relations`

#### 2. 社区模块 (Community)
- **功能**: 帖子发布、浏览、评论、点赞、标签分类
- **API**: `/v1/community/*`
- **核心表**: `posts`, `comments`, `tags`, `user_interactions`

#### 3. 虚拟形象模块 (Avatar)
- **功能**: 虚拟形象管理、道具系统、聊天互动
- **API**: `/v1/avatar/*`
- **核心表**: `avatars`, `user_avatars`, `props`, `prop_categories`, `user_props`, `chat_messages`

#### 4. 小纸条模块 (Message)
- **功能**: AI生成小纸条、付费解锁、情感陪伴
- **API**: `/v1/message/*`
- **核心表**: `messages`, `user_unlock_records`

## 数据库设计

### 核心表结构

- `users` - 用户表
- `pets` - 宠物表
- `posts` - 帖子表
- `comments` - 评论表
- `avatars` - 虚拟形象表
- `props` - 道具表
- `messages` - 小纸条表
- `user_relations` - 用户关系表
- `user_interactions` - 用户互动表

详细的数据库设计请参考 `internal/data/sql/1-init-tables.sql`

## 开发指南

### 环境要求

- **Go**: 1.21+
- **MySQL**: 8.0+
- **Redis**: 6.0+
- **Git**: 最新版本

### 快速开始

#### 1. 克隆项目
```bash
git clone https://github.com/Kononoo/pet-angel.git
cd pet-angel
```

#### 2. 安装依赖
```bash
go mod tidy
```

#### 3. 配置数据库
```bash
# 修改 configs/config.yaml 中的数据库配置
# 执行数据库初始化脚本
mysql -u root -p < internal/data/sql/1-init-tables.sql
```

#### 4. 生成代码
```bash
# 生成API代码
make api

# 生成依赖注入代码
make wire
```

#### 5. 运行项目
```bash
# 开发环境运行
make run

# 或者直接运行
go run cmd/pet-angel/main.go cmd/pet-angel/wire_gen.go -conf ./configs
```

### 项目结构

```
pet-angel/
├── api/                    # API定义
│   ├── user/v1/           # 用户API
│   ├── community/v1/      # 社区API
│   ├── avatar/v1/         # 虚拟形象API
│   └── message/v1/        # 小纸条API
├── cmd/pet-angel/         # 应用入口
│   ├── main.go           # 主程序
│   └── wire.go           # 依赖注入配置
├── internal/              # 内部实现
│   ├── biz/              # 业务逻辑层
│   │   ├── user.go       # 用户业务逻辑
│   │   ├── community.go  # 社区业务逻辑
│   │   ├── avatar.go     # 虚拟形象业务逻辑
│   │   ├── message.go    # 小纸条业务逻辑
│   │   ├── errors.go     # 错误定义
│   │   └── provider.go   # 依赖注入Provider
│   ├── data/             # 数据访问层
│   │   ├── user.go       # 用户数据访问
│   │   ├── community.go  # 社区数据访问
│   │   ├── avatar.go     # 虚拟形象数据访问
│   │   ├── message.go    # 小纸条数据访问
│   │   └── sql/          # 数据库脚本
│   ├── service/          # 服务层
│   │   ├── user.go       # 用户服务
│   │   ├── community.go  # 社区服务
│   │   ├── avatar.go     # 虚拟形象服务
│   │   └── message.go    # 小纸条服务
│   ├── server/           # 服务器层
│   │   ├── grpc.go       # gRPC服务器
│   │   └── http.go       # HTTP服务器
│   └── conf/             # 配置管理
├── configs/              # 配置文件
├── third_party/          # 第三方依赖
├── Makefile              # 构建脚本
├── go.mod               # Go模块文件
├── go.sum               # Go依赖校验
└── README.md            # 项目文档
```

### 开发命令

```bash
# 生成API代码
make api

# 生成依赖注入代码
make wire

# 运行项目
make run

# 构建项目
make build

# 运行测试
make test

# 清理构建文件
make clean
```

## API文档

项目使用gRPC + HTTP双协议，API文档详见各模块的proto文件：

- [用户API](api/user/v1/user.proto)
- [社区API](api/community/v1/community.proto)
- [虚拟形象API](api/avatar/v1/avatar.proto)
- [小纸条API](api/message/v1/message.proto)

### API示例

#### 用户注册
```bash
curl -X POST "http://localhost:8000/v1/user/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "nickname": "测试用户",
    "email": "test@example.com"
  }'
```

#### 获取帖子列表
```bash
curl -X GET "http://localhost:8000/v1/community/posts?page=1&page_size=10"
```

## 部署指南

### Docker部署

```bash
# 构建镜像
docker build -t pet-angel .

# 运行容器
docker run -d \
  --name pet-angel \
  -p 8000:8000 \
  -p 9000:9000 \
  -v /path/to/configs:/app/configs \
  pet-angel
```

### 生产环境配置

1. **数据库配置**: 使用生产环境MySQL集群
2. **Redis配置**: 使用生产环境Redis集群
3. **文件存储**: 配置CDN和对象存储
4. **监控告警**: 配置Prometheus + Grafana
5. **日志管理**: 配置ELK日志系统

## 贡献指南

### 开发规范

1. **代码风格**: 遵循Go官方代码规范
2. **注释要求**: 所有公共函数必须有注释
3. **测试覆盖**: 核心功能必须有单元测试
4. **提交规范**: 使用约定式提交格式

### 提交PR流程

1. Fork项目到个人仓库
2. 创建功能分支: `git checkout -b feature/xxx`
3. 提交代码: `git commit -m "feat: add xxx feature"`
4. 推送分支: `git push origin feature/xxx`
5. 创建Pull Request

### 问题反馈

如果您发现任何问题或有改进建议，请通过以下方式反馈：

- [GitHub Issues](https://github.com/Kononoo/pet-angel/issues)
- [GitHub Discussions](https://github.com/Kononoo/pet-angel/discussions)

## 许可证

本项目采用 [MIT License](LICENSE) 许可证。

## 联系我们

- **项目地址**: https://github.com/Kononoo/pet-angel
- **问题反馈**: https://github.com/Kononoo/pet-angel/issues
- **功能讨论**: https://github.com/Kononoo/pet-angel/discussions

---

**Pet Angel** - 让每一个宠物主人都能获得更好的养宠体验 🐾

