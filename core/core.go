package core

import (
	"context"
	"fmt"
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

// WithAutoParsedBody sets the body and automatically parses XML tags
func WithAutoParsedBody(body string) Option {
	return func(skill *schema.Skill) {
		skill.Body = body
		skill.ParseXMLTags() // 自动解析
	}
}

// WithParsedBody is an alias for WithAutoParsedBody for backward compatibility
// 这是 WithAutoParsedBody 的别名，用于向后兼容
func WithParsedBody(body string) Option {
	return WithAutoParsedBody(body)
}

// Embeddable XML tag functions for constructing body content
// 用于构建 Body 内容的可嵌入 XML 标记函数

// EmbedScript generates a <script>name</script> format string
func EmbedScript(name string) string {
	return string(fmt.Sprintf("<script>%s</script>", name))
}

// EmbedReference generates a <reference>name</reference> format string
func EmbedReference(name string) string {
	return string(fmt.Sprintf("<reference>%s</reference>", name))
}

// EmbedAsset generates a <asset>name</asset> format string
func EmbedAsset(name string) string {
	return string(fmt.Sprintf("<asset>%s</asset>", name))
}

// HasXMLTags checks if the skill's body contains XML tags
// 检查 Skill 的 Body 是否包含 XML 标记
func HasXMLTags(skill *schema.Skill) bool {
	return skill.HasXMLTags()
}

// GetScriptNames gets all script names referenced in the skill's body
// 获取 Body 中引用的所有脚本名称
func GetScriptNames(skill *schema.Skill) []string {
	return skill.GetScriptNames()
}

// GetReferenceNames gets all reference names in the skill's body
// 获取 Body 中引用的所有参考文献名称
func GetReferenceNames(skill *schema.Skill) []string {
	return skill.GetReferenceNames()
}

// GetAssetNames gets all asset names in the skill's body
// 获取 Body 中引用的所有资产名称
func GetAssetNames(skill *schema.Skill) []string {
	return skill.GetAssetNames()
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

// WithResourceProvider sets the ResourceProvider for a skill
// This allows the skill to dynamically load resources from external sources
func WithResourceProvider(provider resources.ResourceProvider) Option {
	return func(skill *schema.Skill) {
		skill.Provider = provider
	}
}

// CreateRemoteScript creates a new RemoteScript with the given name and client
func CreateRemoteScript(name string, client resources.RemoteScriptClient) resources.Script {
	return resources.NewRemoteScript(name, client)
}

// CreateInlineProvider creates a new InlineProvider for in-memory resources
func CreateInlineProvider() *resources.InlineProvider {
	return resources.NewInlineProvider()
}

// CreateCompositeProvider creates a new CompositeProvider that combines multiple providers
func CreateCompositeProvider(providers ...resources.ResourceProvider) *resources.CompositeProvider {
	return resources.NewCompositeProvider(providers...)
}

// CreateCachingProvider creates a new CachingProvider that caches resources
func CreateCachingProvider(provider resources.ResourceProvider) *resources.CachingProvider {
	return resources.NewCachingProvider(provider)
}

// CreateLazyLoadingProvider creates a new LazyLoadingProvider that loads resources on demand
func CreateLazyLoadingProvider(loader func(ctx context.Context) (resources.ResourceProvider, error)) *resources.LazyLoadingProvider {
	return resources.NewLazyLoadingProvider(loader)
}
