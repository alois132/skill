package schema

import (
	"context"
	"encoding/json"
	"errors"

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
