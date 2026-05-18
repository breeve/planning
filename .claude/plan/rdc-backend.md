# 实施计划：研发中心 Go 后端

## 任务类型
- [x] 后端

## 技术方案

基于 spec.md 数据模型，使用 Go 标准项目布局实现 RESTful API 服务。

### 技术选型

| 组件 | 选择 | 理由 |
|------|------|------|
| HTTP 框架 | Gin | 成熟、高性能、社区活跃 |
| ORM | GORM | 支持多对多、JSON 字段、迁移 |
| 数据库 | SQLite（开发）/ PostgreSQL（生产） | SQLite 零配置便于开发测试 |
| 测试 | testing + testify | 标准库 + 断言工具 |
| 配置 | envconfig | 轻量环境变量配置 |

### 项目结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # 入口
├── internal/
│   ├── model/                   # 数据模型定义
│   │   ├── template.go          # 模板层实体
│   │   ├── version.go           # 版本层实体
│   │   └── relations.go         # 多对多关联表
│   ├── repository/              # 数据访问层
│   │   ├── product.go
│   │   ├── subsystem.go
│   │   ├── application.go
│   │   ├── component.go
│   │   ├── artifact.go
│   │   ├── component_version.go
│   │   ├── product_version.go
│   │   └── delivery_plan.go
│   ├── service/                 # 业务逻辑层
│   │   ├── product.go
│   │   ├── subsystem.go
│   │   ├── application.go
│   │   ├── component.go
│   │   ├── artifact.go
│   │   ├── component_version.go
│   │   ├── product_version.go
│   │   └── delivery_plan.go
│   ├── handler/                 # HTTP 处理器
│   │   ├── product.go
│   │   ├── subsystem.go
│   │   ├── application.go
│   │   ├── component.go
│   │   ├── artifact.go
│   │   ├── component_version.go
│   │   ├── product_version.go
│   │   ├── delivery_plan.go
│   │   └── router.go           # 路由注册
│   ├── dto/                     # 请求/响应结构
│   │   ├── request.go
│   │   └── response.go
│   └── database/
│       └── database.go          # DB 初始化与迁移
├── go.mod
├── go.sum
└── Makefile
```

## API 设计

### 模板层 CRUD

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/products | 列表 |
| POST | /api/v1/products | 创建 |
| GET | /api/v1/products/:id | 详情 |
| PUT | /api/v1/products/:id | 更新 |
| DELETE | /api/v1/products/:id | 删除 |
| GET | /api/v1/products/:id/subsystems | 产品下的子系统 |

同理适用于 subsystems、applications、components、artifacts、delivery-plans。

### 多对多关系管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/applications/:id/subsystems | 关联子系统 |
| DELETE | /api/v1/applications/:id/subsystems/:sid | 解除关联 |
| POST | /api/v1/components/:id/applications | 关联应用 |
| DELETE | /api/v1/components/:id/applications/:aid | 解除关联 |

### 版本层操作

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/component-versions | 创建组件版本 |
| GET | /api/v1/component-versions | 列表 |
| GET | /api/v1/component-versions/:id | 详情 |
| POST | /api/v1/product-versions | 创建产品版本 |
| GET | /api/v1/product-versions | 列表 |
| GET | /api/v1/product-versions/:id | 详情 |
| PUT | /api/v1/product-versions/:id/status | 状态流转 |

### 交付方案

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/delivery-plans/:id/product-versions | 添加产品版本 |
| DELETE | /api/v1/delivery-plans/:id/product-versions/:pvid | 移除产品版本 |

### 统一响应格式

```json
{
  "success": true,
  "data": {},
  "error": "",
  "meta": { "total": 0, "page": 1, "limit": 20 }
}
```

## 实施步骤

### 阶段 1：项目初始化与数据层

1. **初始化 Go module** - `go mod init`、安装依赖
2. **定义数据模型** - model 包，包含 GORM tags 和 JSON tags
3. **数据库初始化** - AutoMigrate、连接配置
4. **Repository 层** - 每个实体的 CRUD 方法

预期产物：可运行的数据库迁移 + repository 单元测试

### 阶段 2：业务逻辑层

5. **Service 层** - 封装业务规则
6. **版本创建逻辑** - ComponentVersion 快照组件信息、ProductVersion 快照产品树
7. **状态机** - ProductVersion 状态流转校验（draft→testing→tested→released）
8. **多对多管理** - 关联/解关联逻辑

预期产物：service 单元测试通过

### 阶段 3：HTTP 接口层

9. **Handler 层** - 请求解析、参数校验、调用 service
10. **路由注册** - RESTful 路由组织
11. **统一响应** - 响应封装、错误处理中间件
12. **入口 main.go** - 组装依赖、启动服务

预期产物：API 可通过 curl 测试

### 阶段 4：测试完善

13. **Repository 测试** - 使用 SQLite 内存库
14. **Service 测试** - mock repository
15. **Handler 测试** - httptest 集成测试
16. **Makefile** - build、test、run 命令

预期产物：`make test` 全部通过，覆盖率 ≥ 80%

## 关键文件

| 文件 | 操作 | 说明 |
|------|------|------|
| backend/cmd/server/main.go | 新增 | 服务入口 |
| backend/internal/model/template.go | 新增 | 模板层实体定义 |
| backend/internal/model/version.go | 新增 | 版本层实体定义 |
| backend/internal/model/relations.go | 新增 | 多对多关联表 |
| backend/internal/database/database.go | 新增 | DB 初始化 |
| backend/internal/repository/*.go | 新增 | 数据访问层 |
| backend/internal/service/*.go | 新增 | 业务逻辑层 |
| backend/internal/handler/*.go | 新增 | HTTP 处理器 |
| backend/internal/dto/*.go | 新增 | 请求响应结构 |
| backend/go.mod | 新增 | Go module 定义 |
| backend/Makefile | 新增 | 构建脚本 |

## 核心伪代码

### 数据模型

```go
// model/template.go
type Product struct {
    ID          uint   `gorm:"primaryKey"`
    Name        string `gorm:"uniqueIndex;not null"`
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Subsystem struct {
    ID          uint   `gorm:"primaryKey"`
    Name        string `gorm:"not null"`
    Description string
    ProductID   uint   `gorm:"not null"`
    Product     Product
}

type Application struct {
    ID          uint   `gorm:"primaryKey"`
    Name        string `gorm:"not null"`
    Description string
    Subsystems  []Subsystem `gorm:"many2many:application_subsystems"`
}

type Component struct {
    ID           uint   `gorm:"primaryKey"`
    Name         string `gorm:"not null"`
    Description  string
    Type         string `gorm:"not null;default:'helm'"`
    RepoName     string
    RepoBranch   string
    RepoUser     string
    RepoPasswd   string
    Applications []Application `gorm:"many2many:component_applications"`
}

type Artifact struct {
    ID          uint   `gorm:"primaryKey"`
    Name        string `gorm:"not null"`
    Description string
    ComponentID uint   `gorm:"not null"`
    Component   Component
    Version     string `gorm:"not null"`
    BuiltAt     time.Time
    Registry    string
    RepoName    string
    RepoBranch  string
    RepoCommit  string
}
```

```go
// model/version.go
type ComponentVersion struct {
    ID          uint            `gorm:"primaryKey"`
    ComponentID uint            `gorm:"not null"`
    Component   Component
    Version     string          `gorm:"not null"`
    ArtifactID  uint            `gorm:"not null"`
    Artifact    Artifact
    Snapshot    datatypes.JSON  `gorm:"type:json"`
    Status      string          `gorm:"not null;default:'built'"`
    CreatedBy   string
    CreatedAt   time.Time
}

type ProductVersion struct {
    ID                uint            `gorm:"primaryKey"`
    ProductID         uint            `gorm:"not null"`
    Product           Product
    Version           string          `gorm:"not null"`
    Status            string          `gorm:"not null;default:'draft'"`
    TreeSnapshot      datatypes.JSON  `gorm:"type:json"`
    ComponentVersions []ComponentVersion `gorm:"many2many:product_version_components"`
    CreatedBy         string
    CreatedAt         time.Time
}

type DeliveryPlan struct {
    ID              uint   `gorm:"primaryKey"`
    Name            string `gorm:"not null"`
    Description     string
    ProductVersions []ProductVersion `gorm:"many2many:delivery_plan_product_versions"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### 版本创建逻辑

```go
// service/component_version.go
func (s *ComponentVersionService) Create(componentID, artifactID uint, version, createdBy string) (*model.ComponentVersion, error) {
    component, err := s.componentRepo.FindByID(componentID)
    // snapshot component info
    snapshot := map[string]interface{}{
        "name": component.Name, "type": component.Type,
        "repo_name": component.RepoName, "repo_branch": component.RepoBranch,
    }
    snapshotJSON, _ := json.Marshal(snapshot)
    cv := &model.ComponentVersion{
        ComponentID: componentID, ArtifactID: artifactID,
        Version: version, Snapshot: snapshotJSON,
        Status: "built", CreatedBy: createdBy,
    }
    return s.repo.Create(cv)
}

// service/product_version.go
func (s *ProductVersionService) Create(productID uint, version string, componentVersionIDs []uint, createdBy string) (*model.ProductVersion, error) {
    product, _ := s.productRepo.FindByID(productID)
    subsystems, _ := s.subsystemRepo.FindByProductID(productID)
    // build tree snapshot
    tree := buildTreeSnapshot(product, subsystems, applications, components)
    treeJSON, _ := json.Marshal(tree)
    pv := &model.ProductVersion{
        ProductID: productID, Version: version,
        Status: "draft", TreeSnapshot: treeJSON,
        CreatedBy: createdBy,
    }
    // associate component versions
    return s.repo.CreateWithComponentVersions(pv, componentVersionIDs)
}
```

### 状态机

```go
// service/product_version.go
var validTransitions = map[string][]string{
    "draft":   {"testing"},
    "testing": {"tested"},
    "tested":  {"released"},
}

func (s *ProductVersionService) UpdateStatus(id uint, newStatus string) error {
    pv, _ := s.repo.FindByID(id)
    allowed := validTransitions[pv.Status]
    if !contains(allowed, newStatus) {
        return fmt.Errorf("invalid transition: %s → %s", pv.Status, newStatus)
    }
    return s.repo.UpdateStatus(id, newStatus)
}
```

## 风险与缓解

| 风险 | 缓解措施 |
|------|----------|
| 多对多关系复杂，GORM 预加载性能 | 按需 Preload，避免 N+1 |
| tree_snapshot JSON 结构变更 | 定义明确的 snapshot struct，版本化 |
| 状态流转并发冲突 | 使用乐观锁（GORM 的 version 字段）|
| repo_passwd 明文存储 | 当前阶段暂不处理，后续引入加密 |

## 测试策略

- **Repository 层**：SQLite 内存库，测试实际 SQL 行为
- **Service 层**：interface mock（testify/mock），测试业务逻辑
- **Handler 层**：httptest + gin test mode，测试 HTTP 层集成
- **覆盖重点**：版本创建快照逻辑、状态机流转、多对多关联操作

## SESSION_ID（供 /ccg:execute 使用）
- CODEX_SESSION: N/A（CLI 不可用）
- GEMINI_SESSION: N/A（CLI 不可用）
