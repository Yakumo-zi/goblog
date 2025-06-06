# 📁 项目结构说明

> 本文档描述了博客后端项目的目录结构和文件组织

## 🏗️ 整体架构

项目采用 **Clean Architecture** 设计，遵循Go社区最佳实践的标准项目布局。

## 📂 目录结构

```
goblog/                          # 项目根目录
├── .gitignore                   # Git忽略文件配置
├── README.md                    # 项目说明文档
├── PROJECT_STRUCTURE.md         # 项目结构说明（本文件）
├── Makefile                     # 构建和开发工具
├── go.mod                       # Go模块依赖
├── go.sum                       # Go模块校验
│
├── cmd/                         # 应用程序入口
│   └── server/
│       └── main.go              # 主程序入口
│
├── internal/                    # 私有应用代码
│   ├── config/                  # 配置管理
│   │   └── config.go
│   ├── domain/                  # 领域层（核心业务）
│   │   ├── models.go            # 领域模型
│   │   ├── interfaces.go        # 接口定义
│   │   └── errors.go            # 错误定义
│   ├── handler/                 # HTTP处理层
│   │   ├── article.go           # 文章处理器
│   │   ├── category.go          # 分类处理器
│   │   └── tag.go               # 标签处理器
│   ├── service/                 # 业务逻辑层
│   │   ├── article.go           # 文章服务
│   │   ├── auth.go              # 认证服务
│   │   ├── category.go          # 分类服务
│   │   └── tag.go               # 标签服务
│   ├── repository/              # 数据访问层
│   │   ├── article.go           # 文章仓储
│   │   ├── category.go          # 分类仓储
│   │   └── tag.go               # 标签仓储
│   ├── middleware/              # 中间件
│   │   └── auth.go              # JWT认证中间件
│   └── pkg/                     # 内部通用包
│       ├── logger/              # 日志工具
│       │   └── logger.go
│       └── response/            # 响应工具
│           └── response.go
│
├── test/                        # 测试文件
│   ├── service_test.go          # 旧的通用测试（已重构）
│   ├── category_service_test.go # 分类服务测试
│   └── tag_service_test.go      # 标签服务测试
│
├── scripts/                     # 脚本文件
│   └── test_api.sh              # API功能测试脚本
│
├── bin/                         # 编译输出
│   └── goblog                   # 可执行文件
│
└── ent/                         # Ent ORM生成的文件
    ├── schema/                  # 数据库模式定义
    └── ...                      # 自动生成的ORM代码
```

## 🎯 架构层次

### 1. **表示层 (Presentation Layer)**
- `cmd/server/main.go` - 应用程序入口和依赖注入
- `internal/handler/` - HTTP请求处理器

### 2. **业务逻辑层 (Business Logic Layer)**
- `internal/service/` - 业务服务，包含核心业务逻辑
- `internal/domain/` - 领域模型和接口定义

### 3. **数据访问层 (Data Access Layer)**
- `internal/repository/` - 数据仓储，封装数据库操作
- `ent/` - ORM生成的数据访问代码

### 4. **基础设施层 (Infrastructure Layer)**
- `internal/config/` - 配置管理
- `internal/middleware/` - 中间件
- `internal/pkg/` - 通用工具包

## 📋 设计原则

### ✅ 依赖倒置
- 高层模块不依赖低层模块，都依赖抽象
- 通过接口定义依赖关系

### ✅ 单一职责
- 每个包和模块都有明确的单一职责
- Handler负责HTTP处理，Service负责业务逻辑，Repository负责数据访问

### ✅ 开闭原则
- 通过接口抽象，便于扩展新功能
- 不修改现有代码的情况下添加新特性

### ✅ 依赖注入
- 通过构造函数注入依赖
- 便于单元测试和模块替换

## 🧪 测试策略

### 单元测试
- `test/` 目录包含所有单元测试
- 使用Mock对象隔离测试
- 不需要启动服务器即可运行

### 集成测试
- `scripts/test_api.sh` 提供完整的API功能测试
- 测试真实的HTTP接口

## 🛠️ 开发工具

### Makefile命令
- `make build` - 构建应用程序
- `make test` - 运行单元测试
- `make test-coverage` - 运行测试并生成覆盖率报告
- `make run` - 启动开发服务器
- `make clean` - 清理构建文件
- `make fmt` - 格式化代码

### Git忽略
- `.gitignore` 配置忽略构建产物、日志文件、数据库文件等

## 🔧 项目清理说明

以下文件和目录已在项目整理过程中被删除：

### ❌ 已删除的旧文件
- `main.go` - 旧的主程序文件（已移至 cmd/server/main.go）
- `handlers/` - 旧的处理器目录（已移至 internal/handler/）
- `middleware/` - 旧的中间件目录（已移至 internal/middleware/）
- `pkg/` - 空的包目录（已重新组织至 internal/pkg/）
- `docs/` - 空的文档目录
- `README_NEW.md` - 临时文档（已重命名为 README.md）
- `coverage.html`, `coverage.out` - 测试覆盖率文件（现在被.gitignore忽略）
- `blog.db` - 开发数据库文件（现在被.gitignore忽略）
- `goblog` - 重复的二进制文件（正确位置在 bin/goblog）

### ✅ 清理效果
- 项目结构更加清晰和标准化
- 遵循Go社区最佳实践
- 便于维护和扩展
- 所有功能正常工作（17个单元测试全部通过）

---

> 📝 此项目结构支持企业级开发需求，具有良好的可维护性、可测试性和可扩展性。 