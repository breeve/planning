# 审查任务：验证 harness 内容是否符合用户需求

## 用户原始需求

1. **中间文件规则**：所有Agent在工作的时候，其输出如果会产生中间文件，都需要放到 `.ai/tmp/` 目录下

2. **输出目录规则**：
   - 只有 Docs Agent 产生最终输出，记录到 `docs/ai/` 目录
   - 其他 Agent 不产生最终输出（只有中间输出）

3. **规则名称**：应该是通用的 Agent 规则（agent-rules），不是项目规则（project-rules）

## 需要审查的文件

1. `.ai/harness/rules/agent-rules/SKILL.md` - 规则定义
2. `.ai/harness/agent/docs/agent.md` - Docs Agent
3. `.ai/harness/agent/reviewer/agent.md` - Reviewer Agent
4. `.ai/harness/README.md` - 说明文档

## 审查标准

请检查：
- 规则文件是否包含正确的3条规则
- Docs Agent 的输出目录是否为 `docs/ai/`
- 其他 Agent 的输出目录是否为 `.ai/tmp/`
- 规则是否为通用规则（不是项目特定）
- README 是否正确说明