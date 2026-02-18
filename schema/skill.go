package schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alois132/skill/schema/resources"
	"github.com/alois132/skill/util"
)

type Skill struct {
	Metadata   *SkillMetadata         `json:"metadata"`
	Body       string                 `json:"body"`
	Scripts    []resources.Script     `json:"scripts"`
	References []*resources.Reference `json:"references"`
	Assets     []*resources.Asset     `json:"assets"`

	// 内部缓存字段（不参与序列化）
	parsedTags []util.XMLTag `json:"-"`
	parsed     bool          `json:"-"`
}

type SkillMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (skill *Skill) Glance() (metadata string) {
	m, _ := json.Marshal(skill.Metadata)
	return string(m)
}

func (skill *Skill) Inspect() (body string) {
	return skill.Body
}

func (skill *Skill) UseScript(ctx context.Context, name string, args string) (result string, err error) {
	// 遍历scripts查找匹配名称的脚本
	for _, script := range skill.Scripts {
		if script.GetName() == name {
			return script.Run(ctx, args)
		}
	}
	return "", errors.New("script not found: " + name)
}

func (skill *Skill) ReadReference(name string) (string, error) {
	// 遍历references查找匹配名称的参考文献
	for _, ref := range skill.References {
		if ref.Name == name {
			return ref.Body, nil
		}
	}
	return "", errors.New("reference not found: " + name)
}

// ParseXMLTags 解析 Body 中的 XML 标记并缓存
// 支持格式：\u003cscript\u003ename\u003c/script\u003e, \u003creference\u003ename\u003c/reference\u003e, \u003casset\u003ename\u003c/asset\u003e
func (skill *Skill) ParseXMLTags() error {
	tags := util.ParseXMLTags(skill.Body)
	skill.parsedTags = tags
	skill.parsed = true
	return nil
}

// GetParsedTags 获取已解析的 XML 标记
func (skill *Skill) GetParsedTags() []util.XMLTag {
	if !skill.parsed {
		skill.ParseXMLTags()
	}
	return skill.parsedTags
}

// GetScriptNames 获取 Body 中引用的所有脚本名称
func (skill *Skill) GetScriptNames() []string {
	return util.ExtractScriptNames(skill.Body)
}

// GetReferenceNames 获取 Body 中引用的所有参考文献名称
func (skill *Skill) GetReferenceNames() []string {
	return util.ExtractReferenceNames(skill.Body)
}

// GetAssetNames 获取 Body 中引用的所有资产名称
func (skill *Skill) GetAssetNames() []string {
	return util.ExtractAssetNames(skill.Body)
}

// HasXMLTags 检查 Body 中是否包含 XML 标记
func (skill *Skill) HasXMLTags() bool {
	return util.HasXMLTags(skill.Body)
}

// ScriptResult 脚本执行结果
type ScriptResult struct {
	ScriptName string
	Result     string
	Error      error
}

// AutoExecute 自动执行 Body 中所有 \u003cscript\u003e 标记对应的脚本
// 按出现顺序依次执行，返回所有结果
// 如果某个脚本执行失败，会继续执行后续脚本，错误信息会记录在结果中
func (skill *Skill) AutoExecute(ctx context.Context, args string) ([]ScriptResult, error) {
	scriptNames := skill.GetScriptNames()
	if scriptNames == nil {
		return nil, errors.New("no scripts found in body")
	}

	results := make([]ScriptResult, 0, len(scriptNames))

	for _, scriptName := range scriptNames {
		result, err := skill.UseScript(ctx, scriptName, args)
		results = append(results, ScriptResult{
			ScriptName: scriptName,
			Result:     result,
			Error:      err,
		})
	}

	return results, nil
}

// Execute 执行完整的 skill 逻辑
// 1. 解析 Body 中的 XML 标记
// 2. 按顺序执行所有 \u003cscript\u003e 标记对应的脚本
// 3. 返回组合结果字符串
// 这是 Eino 集成的主要入口点
func (skill *Skill) Execute(ctx context.Context, args string) (string, error) {
	results, err := skill.AutoExecute(ctx, args)
	if err != nil {
		return "", err
	}

	// 组合所有结果
	output := fmt.Sprintf("Skill: %s\n", skill.Metadata.Name)
	for i, result := range results {
		output += fmt.Sprintf("\n[%d] Script: %s\n", i+1, result.ScriptName)
		if result.Error != nil {
			output += fmt.Sprintf("Error: %v\n", result.Error)
		} else {
			output += fmt.Sprintf("Result: %s\n", result.Result)
		}
	}

	return output, nil
}
