# Go博客后端 - 重构版本

这是一个经过完全重构的Go博客后端系统，采用现代化的架构设计和最佳实践。

## 🏗️ 架构特性

- **Clean Architecture** - 分层架构，关注点分离
- **依赖注入** - 基于接口的依赖注入，便于测试
- **领域驱动设计** - 明确的领域模型和业务规则
- **标准项目布局** - 遵循Go社区标准项目结构
- **完整测试覆盖** - 单元测试和集成测试
- **统一错误处理** - 标准化的错误响应格式
- **配置管理** - 环境变量配置支持

## 🛠️ 技术栈

- **语言**: Go 1.24+
- **Web框架**: Echo v4
- **ORM**: Ent
- **数据库**: PostgreSQL (支持Docker部署)
- **认证**: JWT
- **测试**: Testify + Mocks
- **日志**: slog (Go标准库)
- **容器化**: Docker + Docker Compose

## 📁 项目结构

```
goblog/
├── cmd/
│   └── server/                 # 应用程序入口
│       └── main.go
├── internal/                   # 私有代码
│   ├── config/                # 配置管理
│   ├── domain/                # 领域模型和接口
│   │   ├── models.go
│   │   ├── interfaces.go
│   │   └── errors.go
│   ├── handler/               # HTTP处理器
│   ├── service/               # 业务逻辑层
│   ├── repository/            # 数据访问层
│   ├── middleware/            # 中间件
│   └── pkg/                   # 内部工具包
│       ├── logger/
│       └── response/
├── test/                      # 测试文件
├── ent/                       # Ent生成的代码
├── Makefile                   # 构建和开发脚本
├── go.mod
└── README.md
```

## 🚀 快速开始

### 方式一：一键启动 (最简单)

```bash
# 1. 克隆项目
git clone <your-repo-url>
cd goblog

# 2. 运行快速启动脚本
./scripts/quick_start.sh
```

这个脚本会自动：
- 检查系统依赖 (Docker, Docker Compose)
- 启动PostgreSQL数据库
- 启动博客后端服务
- 运行数据库迁移
- 测试API连接
- 显示访问信息

访问服务：
- API服务: http://localhost:8080
- pgAdmin管理界面: http://localhost:8081 (admin@goblog.com / admin123)

### 方式二：本地开发

```bash
# 1. 克隆并安装依赖
git clone <your-repo-url>
cd goblog
make deps

# 2. 启动PostgreSQL数据库
docker-compose up -d postgres

# 3. 运行数据库迁移
./scripts/migrate_postgres.sh migrate
```

# 4. 运行测试
make test

# 5. 启动开发服务器
make run
```

### 方式三：Docker Compose (完整部署)

```bash
# 启动完整服务（包括PostgreSQL数据库）
docker-compose up -d

# 检查服务状态
docker-compose ps

# 查看日志
docker-compose logs -f goblog
```

## 🧪 测试策略

### 单元测试

使用Mock对象进行单元测试，测试每个组件的独立功能：

```bash
# 运行单元测试
go test -v ./test/...

# 查看测试覆盖率
make test-coverage
```

示例测试：

```go
func TestArticleService_Create(t *testing.T) {
    // 准备Mock
    mockArticleRepo := new(MockArticleRepository)
    mockCategoryRepo := new(MockCategoryRepository)
    mockTagRepo := new(MockTagRepository)

    // 创建服务
    articleService := service.NewArticleService(mockArticleRepo, mockCategoryRepo, mockTagRepo)

    // 测试逻辑...
}
```

### 集成测试

集成测试使用真实数据库进行端到端测试：

```bash
# 运行集成测试
go test -tags=integration ./test/integration/...
```

## 🔧 开发工具

### Makefile命令

```bash
make help          # 显示所有可用命令
make build         # 构建应用
make run           # 运行服务器
make test          # 运行测试
make test-coverage # 测试覆盖率报告
make fmt           # 格式化代码
make vet           # 代码检查
make clean         # 清理构建文件
make deps          # 安装依赖
make ent-gen       # 生成Ent代码
make dev           # 开发流程
make ci            # CI流程
```

### 环境变量配置

```bash
# 服务器配置
SERVER_PORT=:8080
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s

# 数据库配置
DB_DRIVER=postgres
DB_DSN=host=localhost port=5432 user=goblog password=goblog123 dbname=goblog sslmode=disable

# JWT配置
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# 管理员配置
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123

# 日志配置
LOG_FORMAT=json  # json 或 text
```

## 🗄️ 数据库管理

### PostgreSQL特性

项目已从SQLite迁移到PostgreSQL，获得以下优势：
- **高并发支持** - 更好的并发读写性能
- **ACID完整性** - 完整的事务支持
- **丰富的数据类型** - JSON、数组等高级数据类型支持
- **扩展性** - 支持复制、分片等扩展方案
- **管理工具** - 完整的数据库管理生态

### 数据库迁移工具

使用内置的PostgreSQL迁移脚本：

```bash
# 检查PostgreSQL连接
./scripts/migrate_postgres.sh check

# 创建数据库
./scripts/migrate_postgres.sh create

# 完整迁移（推荐）
./scripts/migrate_postgres.sh migrate

# 测试数据库连接和表
./scripts/migrate_postgres.sh test

# 显示数据库信息
./scripts/migrate_postgres.sh info
```

### 数据库备份和恢复

Docker Compose自动配置了数据库备份：

```bash
# 手动备份
docker exec goblog-postgres pg_dump -U goblog goblog > backup.sql

# 恢复备份
docker exec -i goblog-postgres psql -U goblog goblog < backup.sql

# 查看自动备份
docker exec goblog-postgres-backup ls -la /backups/
```

### pgAdmin Web管理

访问 http://localhost:8081 使用Web界面管理数据库：
- 用户名: admin@goblog.com
- 密码: admin123

连接设置：
- 主机: postgres
- 端口: 5432
- 数据库: goblog
- 用户名: goblog
- 密码: goblog123

```

## 🌐 API文档

### 认证

#### 登录
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 分类API

#### 获取分类列表（公开）
```bash
curl "http://localhost:8080/api/categories"
```

#### 获取单个分类（公开）
```bash
curl "http://localhost:8080/api/categories/1"
```

#### 创建分类（需要认证）
```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"技术","description":"技术相关分类"}'
```

#### 更新分类（需要认证）
```bash
curl -X PUT http://localhost:8080/api/categories/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"技术更新","description":"更新后的描述"}'
```

#### 删除分类（需要认证）
```bash
curl -X DELETE http://localhost:8080/api/categories/1 \
  -H "Authorization: Bearer <token>"
```

### 标签API

#### 获取标签列表（公开）
```bash
curl "http://localhost:8080/api/tags"
```

#### 获取单个标签（公开）
```bash
curl "http://localhost:8080/api/tags/1"
```

#### 创建标签（需要认证）
```bash
curl -X POST http://localhost:8080/api/tags \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Go语言","color":"#00ADD8"}'
```

#### 更新标签（需要认证）
```bash
curl -X PUT http://localhost:8080/api/tags/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Go语言更新","color":"#FF6B35"}'
```

#### 删除标签（需要认证）
```bash
curl -X DELETE http://localhost:8080/api/tags/1 \
  -H "Authorization: Bearer <token>"
```

### 文章API

#### 获取文章列表（分页）
```bash
curl "http://localhost:8080/api/articles?page=1&limit=10&published=true"
```

#### 创建文章（需要认证）
```bash
curl -X POST http://localhost:8080/api/articles \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "文章标题",
    "content": "文章内容",
    "summary": "文章摘要",
    "published": true,
    "category_id": 1,
    "tag_ids": [1, 2]
  }'
```

#### 按分类获取文章
```bash
curl "http://localhost:8080/api/articles/category/1?page=1&limit=10"
```

#### 按标签获取文章
```bash
curl "http://localhost:8080/api/articles/tag/1?page=1&limit=10"
```

#### 备份所有文章（需要认证）
```bash
curl -X GET http://localhost:8080/api/articles/backup \
  -H "Authorization: Bearer <token>" \
  -o "articles_backup.zip"
```

### 🧪 完整API测试

运行完整的API测试脚本：

```bash
# 启动服务器（在一个终端中）
make run

# 在另一个终端中运行测试
./scripts/test_api.sh
```

测试脚本将会：
- ✅ 登录获取JWT token
- ✅ 创建分类和标签
- ✅ 创建文章并关联分类标签
- ✅ 测试所有读取功能（无需认证）
- ✅ 测试更新功能
- ✅ 测试分页和搜索
- ✅ 测试错误处理
- ✅ 测试文章备份功能
- ✅ 清理测试数据

## 🏛️ 架构设计

### 分层架构

1. **Handler层** - HTTP请求处理，参数验证，响应格式化
2. **Service层** - 业务逻辑，流程控制，业务规则验证
3. **Repository层** - 数据访问，数据模型转换
4. **Domain层** - 领域模型，业务接口，错误定义

### 依赖注入

使用接口实现依赖注入，提高代码的可测试性和可维护性：

```go
// 定义接口
type ArticleService interface {
    Create(ctx context.Context, req *ArticleCreateRequest) (*Article, error)
    // ...
}

// 实现注入
func NewArticleHandler(articleService domain.ArticleService) *ArticleHandler {
    return &ArticleHandler{articleService: articleService}
}
```

### 错误处理

统一的错误处理机制：

```go
// 定义领域错误
var (
    ErrNotFound          = errors.New("resource not found")
    ErrInvalidInput      = errors.New("invalid input")
    ErrDuplicateResource = errors.New("resource already exists")
)

// 统一错误响应
func (h *Handler) handleError(c echo.Context, err error) error {
    if errors.Is(err, domain.ErrNotFound) {
        return response.NotFound(c, "资源不存在")
    }
    // ...
}
```

## 🔒 安全特性

- **JWT认证** - 所有写操作需要JWT token
- **参数验证** - 使用validator库进行请求参数验证
- **CORS支持** - 跨域请求支持
- **密码加密** - 使用bcrypt进行密码哈希

## 📊 监控和日志

- **结构化日志** - 使用slog进行结构化日志记录
- **请求日志** - Echo中间件自动记录所有HTTP请求
- **错误跟踪** - 详细的错误信息和堆栈跟踪

## 🚀 部署

### Docker部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/goblog .
EXPOSE 8080
CMD ["./goblog"]
```

### 直接部署

```bash
# 构建
make build

# 运行
./bin/goblog
```

## 🧪 测试示例

运行测试查看重构后的效果：

```bash
# 运行所有测试
make test

# 查看覆盖率
make test-coverage

# 测试特定功能
go test -v -run TestArticleService_Create ./test/
```

## 📈 性能优化

- **数据库连接池** - Ent自动管理数据库连接
- **分页查询** - 避免大量数据查询
- **索引优化** - 数据库字段索引
- **缓存策略** - 可扩展缓存层

## 🤝 贡献指南

1. Fork项目
2. 创建feature分支
3. 提交更改
4. 运行测试 `make ci`
5. 提交Pull Request

## 📝 许可证

MIT License

---

## 🔄 重构改进点

相比原版本，重构后的版本具有以下改进：

✅ **更好的可测试性** - 依赖注入和Mock测试  
✅ **更清晰的架构** - 分层设计和职责分离  
✅ **更好的可维护性** - 标准项目布局和代码组织  
✅ **更强的类型安全** - 接口定义和错误处理  
✅ **更完整的测试** - 单元测试和集成测试  
✅ **更好的开发体验** - Makefile和自动化工具  
✅ **更标准的实践** - 遵循Go社区最佳实践

## 🐳 Docker部署

### 构建和运行

```bash
# 构建Docker镜像
make docker-build

# 使用Docker Compose运行
make docker-run

# 查看日志
make docker-logs

# 停止服务
make docker-stop

# 清理资源
make docker-clean
```

### Docker镜像优化特性

- **多阶段构建** - 减小最终镜像大小
- **Scratch基础镜像** - 最小化安全攻击面
- **静态链接** - 无运行时依赖
- **压缩二进制** - 去除调试信息
- **健康检查** - 容器状态监控
- **非root用户** - 提高安全性

### 镜像大小对比

- 传统Go镜像：~300MB+
- 优化后镜像：~10MB（压缩后约3-5MB）

## 💾 备份功能

### 文章备份接口

```bash
# 下载所有文章的备份ZIP文件
curl -X GET http://localhost:8080/api/articles/backup \
  -H "Authorization: Bearer <your-jwt-token>" \
  -o "articles_backup_$(date +%Y%m%d_%H%M%S).zip"
```

### 备份文件内容

- `articles_backup.json` - 完整的文章数据（JSON格式）
- `articles/` - 每个文章的单独文件
- `backup_info.txt` - 备份信息说明

### 使用Makefile测试备份

```bash
# 设置JWT token环境变量
export TOKEN="your-jwt-token-here"

# 测试备份功能
make backup-test
``` 