package resources

import (
	"context"
	"errors"
)

// ResourceProvider 统一资源提供者接口
// 用于从各种来源（内存、文件、远程服务）获取 Skill 的资源
type ResourceProvider interface {
	// GetScript 获取指定名称的脚本
	GetScript(ctx context.Context, name string) (Script, error)
	// GetReference 获取指定名称的参考文档内容
	GetReference(ctx context.Context, name string) (string, error)
	// GetAsset 获取指定名称的资源文件
	GetAsset(ctx context.Context, name string) (*Asset, error)

	// ListScripts 列出所有可用的脚本名称
	ListScripts(ctx context.Context) ([]string, error)
	// ListReferences 列出所有可用的参考文档名称
	ListReferences(ctx context.Context) ([]string, error)
	// ListAssets 列出所有可用的资源文件名称
	ListAssets(ctx context.Context) ([]string, error)
}

// InlineProvider 内联资源提供者
// 从内存中的 Scripts、References、Assets 切片提供资源
type InlineProvider struct {
	Scripts    []Script
	References []*Reference
	Assets     []*Asset
}

// NewInlineProvider 创建一个新的内联资源提供者
func NewInlineProvider() *InlineProvider {
	return &InlineProvider{
		Scripts:    make([]Script, 0),
		References: make([]*Reference, 0),
		Assets:     make([]*Asset, 0),
	}
}

// GetScript 从内存中获取脚本
func (p *InlineProvider) GetScript(ctx context.Context, name string) (Script, error) {
	for _, script := range p.Scripts {
		if script.GetName() == name {
			return script, nil
		}
	}
	return nil, errors.New("script not found: " + name)
}

// GetReference 从内存中获取参考文档
func (p *InlineProvider) GetReference(ctx context.Context, name string) (string, error) {
	for _, ref := range p.References {
		if ref.Name == name {
			return ref.Body, nil
		}
	}
	return "", errors.New("reference not found: " + name)
}

// GetAsset 从内存中获取资源文件
func (p *InlineProvider) GetAsset(ctx context.Context, name string) (*Asset, error) {
	for _, asset := range p.Assets {
		if asset.Name == name {
			return asset, nil
		}
	}
	return nil, errors.New("asset not found: " + name)
}

// ListScripts 列出所有脚本名称
func (p *InlineProvider) ListScripts(ctx context.Context) ([]string, error) {
	names := make([]string, len(p.Scripts))
	for i, script := range p.Scripts {
		names[i] = script.GetName()
	}
	return names, nil
}

// ListReferences 列出所有参考文档名称
func (p *InlineProvider) ListReferences(ctx context.Context) ([]string, error) {
	names := make([]string, len(p.References))
	for i, ref := range p.References {
		names[i] = ref.Name
	}
	return names, nil
}

// ListAssets 列出所有资源文件名称
func (p *InlineProvider) ListAssets(ctx context.Context) ([]string, error) {
	names := make([]string, len(p.Assets))
	for i, asset := range p.Assets {
		names[i] = asset.Name
	}
	return names, nil
}

// AddScript 添加脚本到提供者
func (p *InlineProvider) AddScript(script Script) {
	p.Scripts = append(p.Scripts, script)
}

// AddReference 添加参考文档到提供者
func (p *InlineProvider) AddReference(ref *Reference) {
	p.References = append(p.References, ref)
}

// AddAsset 添加资源文件到提供者
func (p *InlineProvider) AddAsset(asset *Asset) {
	p.Assets = append(p.Assets, asset)
}

// Ensure InlineProvider implements ResourceProvider
var _ ResourceProvider = (*InlineProvider)(nil)
