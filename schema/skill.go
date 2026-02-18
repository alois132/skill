package schema

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/alois132/skill/schema/resources"
)

type Skill struct {
	Metadata   *SkillMetadata         `json:"metadata"`
	Body       string                 `json:"body"`
	Scripts    []resources.Script     `json:"scripts"`
	References []*resources.Reference `json:"references"`
	Assets     []*resources.Asset     `json:"assets"`
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
