# Docs Agent 和相关技能审查

请审查以下生成的文件：

## 1. Agent描述文件
- `.ai/docs/agent.md` - 项目文档处理Agent的描述

## 2. Skill文件
- `.ai/skills/docs-reader/SKILL.md` - 文档读取技能
- `.ai/skills/document-parser/SKILL.md` - 文档解析技能
- `.ai/skills/markdown-converter/SKILL.md` - Markdown转换技能

## 审查标准

### Agent描述审查
1. **完整性**：是否包含中文名、英文名、提示词、调用时机
2. **提示词质量**：是否清晰定义核心能力、工作流程、支持格式
3. **格式规范**：是否符合角色描述模板

### Skill描述审查（参考SKILL.md格式规范）
1. **YAML frontmatter**：
   - name：使用kebab-case
   - description：描述技能功能和触发时机
2. **内容质量**：
   - 清晰的使用说明
   - 支持的文件类型
   - 错误处理
3. **格式一致性**：与Anthropic SKILL.md格式兼容

### 整体一致性
- Agent和Skill之间的功能对应
- 命名规范一致性
- 工作流程衔接

请提供详细的审查结果和改进建议。