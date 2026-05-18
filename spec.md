# 研发中心

研发中心是 DevOps 体系中管理产品全生命周期的子系统，覆盖从产品定义、研发构建、质量验收到版本发布的完整流程。

## 目标

- 建立产品结构的统一描述模型，使产品组成可追溯、可版本化
- 规范产研协作流程，明确各角色在每个阶段的职责边界
- 支撑交付场景的灵活组合，同一产品可按需组装为不同交付方案

## 产品结构

产品采用五层层级模型，自顶向下逐层细化：

```txt
交付方案 → 产品 → 子系统 → 应用 → 组件
```

| 层级 | 定义 | 示例 |
| ------ | ------ | ------ |
| 交付方案 | 面向客户的可交付单元，由一组产品版本组合而成 | 企业版、标准版 |
| 产品 | 独立的业务能力集合，具备独立版本 | 工单系统、知识库 |
| 子系统 | 产品内按业务域划分的模块 | 工单-派单引擎、工单-SLA 管理 |
| 应用 | 可独立部署的技术单元 | order-service、order-web |
| 组件 | 应用内可复用的功能模块 | 通知 SDK、权限拦截器 |

## 角色与职责

| 角色 | 核心职责 | 关键动作 |
| ------ | ------ | ------ |
| 产品 | 定义产品结构与能力，管理发布 | 创建产品、定义版本计划、审批发布 |
| 研发 | 实现产品能力，交付可验证的构建制品 | 补充技术定义、构建制品、提测 |
| 测试 | 验证产品质量，确认版本达到发布标准 | 执行验收、标记测试结论 |
| 交付 | 面向客户场景组装交付方案 | 选择产品版本、定义方案组合 |

## 核心流程

```txt
定义 → 构建 → 提测 → 验收 → 发布 → 交付组装
```

1. **定义**（产品）：创建或迭代产品结构，明确本版本的能力范围与组成
2. **构建**（研发）：基于产品定义实现功能，补充技术层细节（应用、组件），产出构建制品
3. **提测**（研发）：制品就绪后提交测试验收，触发版本状态流转
4. **验收**（测试）：执行测试计划，通过后标记版本为"已验收"
5. **发布**（产品）：确认版本可对外发布，锁定版本内容
6. **交付组装**（交付）：将已发布的产品版本组合为面向客户的交付方案

## 数据模型

数据模型分为两层：

- **模板层**：产品结构的当前定义，可随时修改
- **版本层**：某一时刻的不可变快照，由流程事件触发创建

### 模板层（可变）

#### 交付方案(DeliveryPlan)

| 字段 | 说明 |
| - | - |
| name | 名称 |
| description | 描述 |
| product_versions | 包含的产品版本列表 |

#### 产品(Product)

| 字段 | 说明 |
| - | - |
| name | 名称 |
| description | 描述 |

#### 子系统(Subsystem)

| 字段 | 说明 |
| - | - |
| name | 名称 |
| description | 描述 |
| product | 所属产品 |

#### 应用(Application)

| 字段 | 说明 |
| - | - |
| name | 名称 |
| description | 描述 |
| subsystems | 所属子系统，可以属于多个子系统 |

#### 组件(Component)

| 字段 | 说明 |
| - | - |
| name | 名称 |
| description | 描述 |
| applications | 所属应用，可以属于多个应用 |
| type | 组件类型：helm |
| repo_name | 开发仓库 |
| repo_branch | 开发发布分支 |
| repo_user | 开发仓库用户 |
| repo_passwd | 开发仓库密码 |

#### 制品(Artifact)

| 字段 | 说明 |
| - | - |
| name | 名称 |
| description | 描述 |
| component | 所属组件 |
| version | 版本号 |
| built_at | 构建时间 |
| registry | 制品存储地址（如 OCI registry URL） |
| repo_name | 构建时仓库名（快照） |
| repo_branch | 构建时分支（快照） |
| repo_commit | 构建时提交 ID（快照） |

### 版本层（不可变快照）

版本由流程事件触发创建（提测、发布等），创建后内容不可修改。

#### 组件版本(ComponentVersion)

原子版本单元，绑定一个具体制品。

| 字段 | 说明 |
| - | - |
| component_id | 指向组件模板 |
| version | 版本号 |
| artifact_id | 绑定的制品 |
| snapshot | 创建时组件信息快照（name、type、repo 等） |
| status | built / integrated |
| created_by | 创建人 |
| created_at | 创建时间 |

#### 产品版本(ProductVersion)

聚合版本，记录整棵产品树在某一时刻的完整状态。

| 字段 | 说明 |
| - | - |
| product_id | 指向产品模板 |
| version | 版本号 |
| status | draft / testing / tested / released |
| tree_snapshot | 产品结构快照（产品/子系统/应用的 name、description 等） |
| component_versions | 包含的 ComponentVersion 列表 |
| created_by | 创建人 |
| created_at | 创建时间 |

### 模型关系

```txt
模板层（可变）：
  Product ──1:N── Subsystem ──M:N── Application ──M:N── Component ──1:N── Artifact

版本层（不可变）：
  ProductVersion ──1:N── ComponentVersion ──1:1── Artifact

交付组装：
  DeliveryPlan ──M:N── ProductVersion
```

### 版本创建规则

1. **组件制品触发**：构建产出新 Artifact → 创建 ComponentVersion
2. **向上聚合**：选择一组 ComponentVersion → 快照当前产品树结构 → 创建 ProductVersion
3. **模板变更触发**：修改子系统/应用等结构信息 → 如需发布，创建新 ProductVersion（tree_snapshot 反映变更）
4. **状态流转**：ProductVersion 的 status 承载工作流（draft → testing → tested → released）

## 核心用例

## 系统架构

## 工作流

## 应用规范

> 描述应用的组成规范，包含环境变量、配置类型、发布方式、发布类型等
