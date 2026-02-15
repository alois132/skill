package core

import (
	"github.com/alois132/skill/schema"
	"github.com/alois132/skill/schema/resources"
)

// 不考虑封锁边界，只是为了便捷和美观
type Option func(skill *schema.Skill)

// create skill

// create reference

func WithReferences(refs []*resources.Reference) Option {
	return func(skill *schema.Skill) {
		skill.References = refs
	}
}

func CreateReference(name string, body string) *resources.Reference {
	return &resources.Reference{
		Name: name,
		Body: body,
	}
}

// create script

func WithScripts(scripts []resources.Script) Option {
	return func(skill *schema.Skill) {
		skill.Scripts = scripts
	}
}

func CreateScript[I, O any](name string, fn resources.ScriptFunc[I, O]) resources.Script {
	return &resources.EasyScript[I, O]{
		Name: name,
		Fn:   fn,
	}
}

// glance skill

// inspect skill

// use skill's script

// read skill's reference
