package resources

import (
	"context"
	"errors"
	"fmt"
)

// CompositeProvider 复合资源提供者
// 可以组合多个提供者，按优先级顺序查找资源
type CompositeProvider struct {
	providers []ResourceProvider
}

// NewCompositeProvider 创建一个新的复合资源提供者
func NewCompositeProvider(providers ...ResourceProvider) *CompositeProvider {
	return &CompositeProvider{
		providers: providers,
	}
}

// AddProvider 添加一个资源提供者
func (p *CompositeProvider) AddProvider(provider ResourceProvider) {
	p.providers = append(p.providers, provider)
}

// GetScript 按优先级从所有提供者中查找脚本
func (p *CompositeProvider) GetScript(ctx context.Context, name string) (Script, error) {
	var lastErr error
	for _, provider := range p.providers {
		script, err := provider.GetScript(ctx, name)
		if err == nil {
			return script, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.New("script not found: " + name)
}

// GetReference 按优先级从所有提供者中查找参考文档
func (p *CompositeProvider) GetReference(ctx context.Context, name string) (string, error) {
	var lastErr error
	for _, provider := range p.providers {
		ref, err := provider.GetReference(ctx, name)
		if err == nil {
			return ref, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return "", lastErr
	}
	return "", errors.New("reference not found: " + name)
}

// GetAsset 按优先级从所有提供者中查找资源文件
func (p *CompositeProvider) GetAsset(ctx context.Context, name string) (*Asset, error) {
	var lastErr error
	for _, provider := range p.providers {
		asset, err := provider.GetAsset(ctx, name)
		if err == nil {
			return asset, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.New("asset not found: " + name)
}

// ListScripts 合并所有提供者的脚本列表
func (p *CompositeProvider) ListScripts(ctx context.Context) ([]string, error) {
	nameSet := make(map[string]struct{})
	for _, provider := range p.providers {
		names, err := provider.ListScripts(ctx)
		if err != nil {
			continue // 跳过出错的提供者
		}
		for _, name := range names {
			nameSet[name] = struct{}{}
		}
	}
	names := make([]string, 0, len(nameSet))
	for name := range nameSet {
		names = append(names, name)
	}
	return names, nil
}

// ListReferences 合并所有提供者的参考文档列表
func (p *CompositeProvider) ListReferences(ctx context.Context) ([]string, error) {
	nameSet := make(map[string]struct{})
	for _, provider := range p.providers {
		names, err := provider.ListReferences(ctx)
		if err != nil {
			continue // 跳过出错的提供者
		}
		for _, name := range names {
			nameSet[name] = struct{}{}
		}
	}
	names := make([]string, 0, len(nameSet))
	for name := range nameSet {
		names = append(names, name)
	}
	return names, nil
}

// ListAssets 合并所有提供者的资源文件列表
func (p *CompositeProvider) ListAssets(ctx context.Context) ([]string, error) {
	nameSet := make(map[string]struct{})
	for _, provider := range p.providers {
		names, err := provider.ListAssets(ctx)
		if err != nil {
			continue // 跳过出错的提供者
		}
		for _, name := range names {
			nameSet[name] = struct{}{}
		}
	}
	names := make([]string, 0, len(nameSet))
	for name := range nameSet {
		names = append(names, name)
	}
	return names, nil
}

// CachingProvider 带缓存的资源提供者装饰器
type CachingProvider struct {
	provider     ResourceProvider
	scriptCache  map[string]Script
	refCache     map[string]string
	assetCache   map[string]*Asset
}

// NewCachingProvider 创建一个新的缓存资源提供者
func NewCachingProvider(provider ResourceProvider) *CachingProvider {
	return &CachingProvider{
		provider:    provider,
		scriptCache: make(map[string]Script),
		refCache:    make(map[string]string),
		assetCache:  make(map[string]*Asset),
	}
}

// GetScript 从缓存或底层提供者获取脚本
func (p *CachingProvider) GetScript(ctx context.Context, name string) (Script, error) {
	if script, ok := p.scriptCache[name]; ok {
		return script, nil
	}
	script, err := p.provider.GetScript(ctx, name)
	if err != nil {
		return nil, err
	}
	p.scriptCache[name] = script
	return script, nil
}

// GetReference 从缓存或底层提供者获取参考文档
func (p *CachingProvider) GetReference(ctx context.Context, name string) (string, error) {
	if ref, ok := p.refCache[name]; ok {
		return ref, nil
	}
	ref, err := p.provider.GetReference(ctx, name)
	if err != nil {
		return "", err
	}
	p.refCache[name] = ref
	return ref, nil
}

// GetAsset 从缓存或底层提供者获取资源文件
func (p *CachingProvider) GetAsset(ctx context.Context, name string) (*Asset, error) {
	if asset, ok := p.assetCache[name]; ok {
		return asset, nil
	}
	asset, err := p.provider.GetAsset(ctx, name)
	if err != nil {
		return nil, err
	}
	p.assetCache[name] = asset
	return asset, nil
}

// ListScripts 列出所有脚本（不缓存）
func (p *CachingProvider) ListScripts(ctx context.Context) ([]string, error) {
	return p.provider.ListScripts(ctx)
}

// ListReferences 列出所有参考文档（不缓存）
func (p *CachingProvider) ListReferences(ctx context.Context) ([]string, error) {
	return p.provider.ListReferences(ctx)
}

// ListAssets 列出所有资源文件（不缓存）
func (p *CachingProvider) ListAssets(ctx context.Context) ([]string, error) {
	return p.provider.ListAssets(ctx)
}

// ClearCache 清除所有缓存
func (p *CachingProvider) ClearCache() {
	p.scriptCache = make(map[string]Script)
	p.refCache = make(map[string]string)
	p.assetCache = make(map[string]*Asset)
}

// ClearScriptCache 清除脚本缓存
func (p *CachingProvider) ClearScriptCache() {
	p.scriptCache = make(map[string]Script)
}

// ClearReferenceCache 清除参考文档缓存
func (p *CachingProvider) ClearReferenceCache() {
	p.refCache = make(map[string]string)
}

// ClearAssetCache 清除资源文件缓存
func (p *CachingProvider) ClearAssetCache() {
	p.assetCache = make(map[string]*Asset)
}

// Ensure CompositeProvider implements ResourceProvider
var _ ResourceProvider = (*CompositeProvider)(nil)

// Ensure CachingProvider implements ResourceProvider
var _ ResourceProvider = (*CachingProvider)(nil)

// LazyLoadingProvider 懒加载资源提供者
// 只在首次访问时从 loader 加载资源
type LazyLoadingProvider struct {
	loader       func(ctx context.Context) (ResourceProvider, error)
	provider     ResourceProvider
	initialized  bool
}

// NewLazyLoadingProvider 创建一个新的懒加载资源提供者
func NewLazyLoadingProvider(loader func(ctx context.Context) (ResourceProvider, error)) *LazyLoadingProvider {
	return &LazyLoadingProvider{
		loader: loader,
	}
}

func (p *LazyLoadingProvider) init(ctx context.Context) error {
	if p.initialized {
		return nil
	}
	provider, err := p.loader(ctx)
	if err != nil {
		return fmt.Errorf("failed to load provider: %w", err)
	}
	p.provider = provider
	p.initialized = true
	return nil
}

// GetScript 获取脚本
func (p *LazyLoadingProvider) GetScript(ctx context.Context, name string) (Script, error) {
	if err := p.init(ctx); err != nil {
		return nil, err
	}
	return p.provider.GetScript(ctx, name)
}

// GetReference 获取参考文档
func (p *LazyLoadingProvider) GetReference(ctx context.Context, name string) (string, error) {
	if err := p.init(ctx); err != nil {
		return "", err
	}
	return p.provider.GetReference(ctx, name)
}

// GetAsset 获取资源文件
func (p *LazyLoadingProvider) GetAsset(ctx context.Context, name string) (*Asset, error) {
	if err := p.init(ctx); err != nil {
		return nil, err
	}
	return p.provider.GetAsset(ctx, name)
}

// ListScripts 列出脚本
func (p *LazyLoadingProvider) ListScripts(ctx context.Context) ([]string, error) {
	if err := p.init(ctx); err != nil {
		return nil, err
	}
	return p.provider.ListScripts(ctx)
}

// ListReferences 列出参考文档
func (p *LazyLoadingProvider) ListReferences(ctx context.Context) ([]string, error) {
	if err := p.init(ctx); err != nil {
		return nil, err
	}
	return p.provider.ListReferences(ctx)
}

// ListAssets 列出资源文件
func (p *LazyLoadingProvider) ListAssets(ctx context.Context) ([]string, error) {
	if err := p.init(ctx); err != nil {
		return nil, err
	}
	return p.provider.ListAssets(ctx)
}

// Ensure LazyLoadingProvider implements ResourceProvider
var _ ResourceProvider = (*LazyLoadingProvider)(nil)
