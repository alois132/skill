package eino

import (
	"github.com/alois132/skill/schema"
	"github.com/cloudwego/eino/components/tool"
)

// ToTools 将多个 Skill 转换为 Eino Tools
// 返回 n+2 个 Tool：n 个 SkillTool + 1 个 UseScriptTool + 1 个 ReadReferenceTool
func ToTools(skills ...*schema.Skill) []tool.BaseTool {
	if len(skills) == 0 {
		return nil
	}

	// 创建 n 个 SkillTool（每个 skill 一个）
	tools := make([]tool.BaseTool, 0, len(skills)+2)
	for _, skill := range skills {
		tools = append(tools, NewSkillTool(skill))
	}

	// 添加共享的 use_script 和 read_reference 工具
	tools = append(tools, NewUseScriptTool(skills...))
	tools = append(tools, NewReadReferenceTool(skills...))

	return tools
}

// ToInvokableTools 将多个 Skill 转换为 Eino InvokableTools
// 返回的 tools 可以直接用于 ToolsNode
func ToInvokableTools(skills ...*schema.Skill) []tool.InvokableTool {
	if len(skills) == 0 {
		return nil
	}

	// 创建 n 个 SkillTool（每个 skill 一个）
	tools := make([]tool.InvokableTool, 0, len(skills)+2)
	for _, skill := range skills {
		tools = append(tools, NewSkillTool(skill))
	}

	// 添加共享的 use_script 和 read_reference 工具
	tools = append(tools, NewUseScriptTool(skills...))
	tools = append(tools, NewReadReferenceTool(skills...))

	return tools
}
