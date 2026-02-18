package core

import (
	"context"
	"github.com/alois132/skill/schema"
	"github.com/alois132/skill/schema/resources"
)

// 不考虑封锁边界，只是为了便捷和美观
type Option func(skill *schema.Skill)

// create skill

// CreateSkill creates a new skill with the given metadata and options
func CreateSkill(name string, description string, opts ...Option) *schema.Skill {
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{
			Name:        name,
			Description: description,
		},
		Body:       "",
		Scripts:    []resources.Script{},
		References: []*resources.Reference{},
		Assets:     []*resources.Asset{},
	}

	// Apply all options
	for _, opt := range opts {
		opt(skill)
	}

	return skill
}

// create reference

// WithReferences adds multiple references to a skill
func WithReferences(refs []*resources.Reference) Option {
	return func(skill *schema.Skill) {
		skill.References = append(skill.References, refs...)
	}
}

// WithReference adds a single reference to a skill
func WithReference(name string, body string) Option {
	return func(skill *schema.Skill) {
		skill.References = append(skill.References, &resources.Reference{
			Name: name,
			Body: body,
		})
	}
}

// CreateReference creates a new reference with the given name and body
func CreateReference(name string, body string) *resources.Reference {
	return &resources.Reference{
		Name: name,
		Body: body,
	}
}

// create script

// WithScripts adds multiple scripts to a skill
func WithScripts(scripts []resources.Script) Option {
	return func(skill *schema.Skill) {
		skill.Scripts = append(skill.Scripts, scripts...)
	}
}

// WithScript adds a single script to a skill
func WithScript(script resources.Script) Option {
	return func(skill *schema.Skill) {
		skill.Scripts = append(skill.Scripts, script)
	}
}

// CreateScript creates a new EasyScript with the given name and function
func CreateScript[I, O any](name string, fn resources.ScriptFunc[I, O]) resources.Script {
	return &resources.EasyScript[I, O]{
		Name: name,
		Fn:   fn,
	}
}

// create asset

// WithAssets adds multiple assets to a skill
func WithAssets(assets []*resources.Asset) Option {
	return func(skill *schema.Skill) {
		skill.Assets = append(skill.Assets, assets...)
	}
}

// WithAsset adds a single asset to a skill
func WithAsset(asset *resources.Asset) Option {
	return func(skill *schema.Skill) {
		skill.Assets = append(skill.Assets, asset)
	}
}

// CreateAsset creates a new asset with the given name, bytes, and extension
func CreateAsset(name string, data []byte, ext resources.AssetExt) *resources.Asset {
	return &resources.Asset{
		Name: name,
		Bytes: data,
		Ext:  ext,
	}
}

// WithBody sets the body of a skill
func WithBody(body string) Option {
	return func(skill *schema.Skill) {
		skill.Body = body
	}
}

// glance skill

// Glance returns a glance view of the skill's metadata
func Glance(skill *schema.Skill) (metadata string) {
	return skill.Glance()
}

// inspect skill

// Inspect returns the detailed body of the skill
func Inspect(skill *schema.Skill) (body string) {
	return skill.Inspect()
}

// use skill's script

// UseScript executes a script by name on the given skill
func UseScript(ctx context.Context, skill *schema.Skill, name string, args string) (result string, err error) {
	return skill.UseScript(ctx, name, args)
}

// read skill's reference

// ReadReference reads a reference by name from the given skill
func ReadReference(skill *schema.Skill, name string) (string, error) {
	return skill.ReadReference(name)
}
