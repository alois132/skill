package schema

import (
	"context"
	"encoding/json"
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

}
