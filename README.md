# 功能

- 研发中心 负责制品构建
- 制品中心 负责管理制品
- 产品中心 负责定义产品
- 配置中心 负责定义监督相关的一些配置元数据
- 规划中心 负责规划产品的交付配置
- 部署中心 负责部署产品到环境

# 研发中心

以组件为纬度进行管理，包含组件的源码仓库 分支 release-commit-id 约束出来组件唯一版本

组件需要包含

- app-life-tools
- install
- upgrade
- rollback
- uninstall

app-config

- provide config
- deps config
- trait config

app-publish-tools

- publish

# 制品中心

纯粹制品库 用于记录研发中心推送的制品，以组件为纬度进行管理

# 产品中心

以"_交付方案-产品-子系统-组件_"为纬度组织产品信息

# 配置中心

以业务scope划分，定义交付相关配置元数据

# 规划中心

提供环境管理，交付规划 两个核心能力

# 部署中心

根据规划输出以及组件物料，完成环境部署
