package eino

import (
	"context"
	"encoding/json"
	"fmt"

	skillschema "github.com/alois132/skill/schema"
	"github.com/cloudwego/eino/components/tool"
	einosch "github.com/cloudwego/eino/schema"
)

// SkillTool 将 Skill 封装为 Eino Tool
// 用于渐进式披露：input 为空，output 为 skill.Body
type SkillTool struct {
	skill *skillschema.Skill
}

// NewSkillTool 创建一个新的 SkillTool
func NewSkillTool(skill *skillschema.Skill) *SkillTool {
	return &SkillTool{skill: skill}
}

// Info 返回 Tool 的元信息
func (t *SkillTool) Info(ctx context.Context) (*einosch.ToolInfo, error) {
	// 空参数表示不需要输入
	return &einosch.ToolInfo{
		Name: t.skill.Metadata.Name,
		Desc: t.skill.Metadata.Description,
	}, nil
}

// InvokableRun 执行 Tool
// input 为空，直接返回 skill.Body
func (t *SkillTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	return t.skill.Body, nil
}

// UseScriptRequest use_script 工具的请求参数
type UseScriptRequest struct {
	SkillName  string `json:"skill_name"`
	ScriptName string `json:"script_name"`
	Args       string `json:"args"`
}

// UseScriptTool 执行 Skill 中的特定脚本
type UseScriptTool struct {
	skills map[string]*skillschema.Skill // skill name -> skill
}

// NewUseScriptTool 创建一个新的 UseScriptTool
func NewUseScriptTool(skills ...*skillschema.Skill) *UseScriptTool {
	skillMap := make(map[string]*skillschema.Skill, len(skills))
	for _, skill := range skills {
		skillMap[skill.Metadata.Name] = skill
	}
	return &UseScriptTool{skills: skillMap}
}

// Info 返回 Tool 的元信息
func (t *UseScriptTool) Info(ctx context.Context) (*einosch.ToolInfo, error) {
	params := map[string]*einosch.ParameterInfo{
		"skill_name": {
			Type: einosch.Object,
			Desc: "The name of the skill to use",
		},
		"script_name": {
			Type: einosch.Object,
			Desc: "The name of the script to execute (found in <script> tags)",
		},
		"args": {
			Type: einosch.Object,
			Desc: "JSON string of arguments to pass to the script",
		},
	}
	// 标记 skill_name 和 script_name 为 required
	params["skill_name"].Required = true
	params["script_name"].Required = true

	info := &einosch.ToolInfo{
		Name: "use_script",
		Desc: "Execute a specific script from a skill. Call this after getting the skill body to run scripts referenced in <script> tags.",
	}
	info.ParamsOneOf = einosch.NewParamsOneOfByParams(params)
	return info, nil
}

// InvokableRun 执行 Tool
func (t *UseScriptTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	var req UseScriptRequest
	if err := json.Unmarshal([]byte(argumentsInJSON), &req); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	skill, ok := t.skills[req.SkillName]
	if !ok {
		return "", fmt.Errorf("skill not found: %s", req.SkillName)
	}

	return skill.UseScript(ctx, req.ScriptName, req.Args)
}

// ReadReferenceRequest read_reference 工具的请求参数
type ReadReferenceRequest struct {
	SkillName     string `json:"skill_name"`
	ReferenceName string `json:"reference_name"`
}

// ReadReferenceTool 读取 Skill 中的参考文献
type ReadReferenceTool struct {
	skills map[string]*skillschema.Skill
}

// NewReadReferenceTool 创建一个新的 ReadReferenceTool
func NewReadReferenceTool(skills ...*skillschema.Skill) *ReadReferenceTool {
	skillMap := make(map[string]*skillschema.Skill, len(skills))
	for _, skill := range skills {
		skillMap[skill.Metadata.Name] = skill
	}
	return &ReadReferenceTool{skills: skillMap}
}

// Info 返回 Tool 的元信息
func (t *ReadReferenceTool) Info(ctx context.Context) (*einosch.ToolInfo, error) {
	params := map[string]*einosch.ParameterInfo{
		"skill_name": {
			Type:        einosch.Object,
			Desc:        "The name of the skill to use",
			Required:    true,
		},
		"reference_name": {
			Type:        einosch.Object,
			Desc:        "The name of the reference to read (found in <reference> tags)",
			Required:    true,
		},
	}

	info := &einosch.ToolInfo{
		Name: "read_reference",
		Desc: "Read a reference document from a skill. Call this after getting the skill body to read references mentioned in <reference> tags.",
	}
	info.ParamsOneOf = einosch.NewParamsOneOfByParams(params)
	return info, nil
}

// InvokableRun 执行 Tool
func (t *ReadReferenceTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	var req ReadReferenceRequest
	if err := json.Unmarshal([]byte(argumentsInJSON), &req); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	skill, ok := t.skills[req.SkillName]
	if !ok {
		return "", fmt.Errorf("skill not found: %s", req.SkillName)
	}

	return skill.ReadReference(req.ReferenceName)
}
