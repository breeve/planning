# Reviewer Agent审查

## 审查类型
DevOps功能分解和架构设计审查

## 审查对象

### 1. 功能分解结果
- 位置：.ai/devops_analysis.md
- 内容：DevOps产品功能分解（14项功能，4个类别）

### 2. Drawio架构图
- 位置：.ai/devops_architecture.drawio.svg
- 内容：DevOps功能架构图

## 本次审查要点

根据格式特定审查规则，重点审查：

### Drawio XML审查
1. 节点位置检查：坐标是否有效
2. 边连接检查：source/target是否正确
3. 布局质量检查：是否有边重叠
4. 布局优化：是否使用水平分支布局

## 审查请求

请按照以下规则进行审查：
- 如果通过，返回approve
- 如果需要修改，返回needs_modification并说明原因
- 如果不合理，返回reject

## 审查输出要求
- review_result
- feedback
- risk_assessment
- approval_status
- next_steps