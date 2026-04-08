---
name: agent-rules
description: Agent通用规则定义，定义Agent工作时需要遵守的规范。触发词：工作规范、输出规则、文件规则。
---

# Agent 通用规则

本文件定义了Agent在工作时应遵守的通用规则。

## 规则1：中间文件规则

**规则**：所有Agent在工作的时候，其输出如果会产生中间文件，都需要放到 `.ai/tmp/` 目录下

**说明**：
- 中间文件是指处理过程中的临时文件，非最终交付物
- 每个Agent在处理任务前，应先创建 `.ai/tmp/` 目录（如果不存在）
- 中间文件包括：审查输入、审查报告、临时分析文件等
- 任务完成后，根据需要决定是否清理中间文件

**示例**：
```
.ai/tmp/reviewer_input.md
.ai/tmp/reviewer_report.md
.ai/tmp/temp_analysis.md
```

## 规则2：输出目录规则

**规则**：只有 Docs Agent 产生最终输出，其他 Agent 不产生最终输出

**说明**：
- **Docs Agent**：最终输出保存到 `docs/ai/` 目录
- **其他 Agent**（search, engineering, system-planner, reviewer）：仅输出中间文件到 `.ai/tmp/`，不产生最终输出
- 其他 Agent 如需输出审查结果或报告，也保存到 `.ai/tmp/`

**示例**：
```
# Docs Agent 最终输出
docs/ai/analysis-report.md
docs/ai/architecture-diagram.svg

# 其他 Agent 中间输出
.ai/tmp/reviewer_input.md
.ai/tmp/analysis_result.md
```

## 规则3：文件命名规范

- 使用英文命名
- 遵循 kebab-case（如 `analysis-report.md`）
- Agent定义文件：`agent.md`
- Skill定义文件：`SKILL.md`
- 规则定义文件：`SKILL.md`