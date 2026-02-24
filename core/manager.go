package core

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/alois132/skill/schema"
	"github.com/alois132/skill/schema/resources"
	"github.com/alois132/skill/schema/store"
)

// SkillManager 统一管理 Skill 的加载、缓存和生命周期
type SkillManager struct {
	store     store.SkillStore
	cache     map[string]*schema.Skill
	mu        sync.RWMutex
	providers map[string]resources.ResourceProvider // skill name -> provider
}

// ManagerOption SkillManager 的配置选项
type ManagerOption func(*SkillManager)

// NewSkillManager 创建一个新的 Skill 管理器
func NewSkillManager(store store.SkillStore, opts ...ManagerOption) *SkillManager {
	m := &SkillManager{
		store:     store,
		cache:     make(map[string]*schema.Skill),
		providers: make(map[string]resources.ResourceProvider),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// WithManagerResourceProvider 为指定的 Skill 设置资源提供者
func WithManagerResourceProvider(skillName string, provider resources.ResourceProvider) ManagerOption {
	return func(m *SkillManager) {
		m.providers[skillName] = provider
	}
}

// GetSkill 获取指定名称的 Skill
// 优先从缓存获取，如果缓存未命中则从 Store 加载
func (m *SkillManager) GetSkill(ctx context.Context, name string) (*schema.Skill, error) {
	// 1. 尝试从缓存获取
	m.mu.RLock()
	if skill, ok := m.cache[name]; ok {
		m.mu.RUnlock()
		return skill, nil
	}
	m.mu.RUnlock()

	// 2. 从 Store 加载
	if m.store == nil {
		return nil, errors.New("skill store not configured")
	}

	skill, err := m.store.Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to load skill from store: %w", err)
	}

	// 3. 如果该 Skill 有配置 ResourceProvider，则设置
	if provider, ok := m.providers[name]; ok {
		skill.Provider = provider
	}

	// 4. 存入缓存
	m.mu.Lock()
	m.cache[name] = skill
	m.mu.Unlock()

	return skill, nil
}

// RegisterSkill 直接注册一个 Skill 到管理器（不经过 Store）
func (m *SkillManager) RegisterSkill(skill *schema.Skill) error {
	if skill == nil {
		return errors.New("skill cannot be nil")
	}
	if skill.Metadata == nil || skill.Metadata.Name == "" {
		return errors.New("skill metadata name cannot be empty")
	}

	name := skill.Metadata.Name

	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache[name] = skill
	return nil
}

// SaveSkill 保存 Skill 到 Store 并更新缓存
func (m *SkillManager) SaveSkill(ctx context.Context, skill *schema.Skill) error {
	if m.store == nil {
		return errors.New("skill store not configured")
	}

	if err := m.store.Put(ctx, skill); err != nil {
		return fmt.Errorf("failed to save skill to store: %w", err)
	}

	// 更新缓存
	m.mu.Lock()
	if skill.Metadata != nil {
		m.cache[skill.Metadata.Name] = skill
	}
	m.mu.Unlock()

	return nil
}

// ReloadSkill 重新从 Store 加载指定的 Skill
func (m *SkillManager) ReloadSkill(ctx context.Context, name string) (*schema.Skill, error) {
	if m.store == nil {
		return nil, errors.New("skill store not configured")
	}

	// 从 Store 重新加载
	skill, err := m.store.Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to reload skill: %w", err)
	}

	// 如果该 Skill 有配置 ResourceProvider，则设置
	if provider, ok := m.providers[name]; ok {
		skill.Provider = provider
	}

	// 更新缓存
	m.mu.Lock()
	m.cache[name] = skill
	m.mu.Unlock()

	return skill, nil
}

// DeleteSkill 从 Store 和缓存中删除指定的 Skill
func (m *SkillManager) DeleteSkill(ctx context.Context, name string) error {
	if m.store == nil {
		return errors.New("skill store not configured")
	}

	if err := m.store.Delete(ctx, name); err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	// 从缓存中移除
	m.mu.Lock()
	delete(m.cache, name)
	m.mu.Unlock()

	return nil
}

// ListSkills 列出所有可用的 Skill 元数据
func (m *SkillManager) ListSkills(ctx context.Context) ([]*schema.SkillMetadata, error) {
	if m.store == nil {
		// 如果没有 Store，返回缓存中的 Skill 元数据
		m.mu.RLock()
		defer m.mu.RUnlock()

		metadatas := make([]*schema.SkillMetadata, 0, len(m.cache))
		for _, skill := range m.cache {
			if skill.Metadata != nil {
				metadatas = append(metadatas, skill.Metadata)
			}
		}
		return metadatas, nil
	}

	return m.store.List(ctx)
}

// UseScript 执行指定 Skill 的脚本
func (m *SkillManager) UseScript(ctx context.Context, skillName string, scriptName string, args string) (string, error) {
	skill, err := m.GetSkill(ctx, skillName)
	if err != nil {
		return "", err
	}

	return skill.UseScript(ctx, scriptName, args)
}

// ReadReference 读取指定 Skill 的参考文档
func (m *SkillManager) ReadReference(ctx context.Context, skillName string, refName string) (string, error) {
	skill, err := m.GetSkill(ctx, skillName)
	if err != nil {
		return "", err
	}

	return skill.ReadReference(refName)
}

// ClearCache 清空 Skill 缓存
func (m *SkillManager) ClearCache() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache = make(map[string]*schema.Skill)
}

// GetCachedSkillNames 获取当前缓存中的所有 Skill 名称
func (m *SkillManager) GetCachedSkillNames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.cache))
	for name := range m.cache {
		names = append(names, name)
	}
	return names
}

// SetResourceProvider 为指定的 Skill 设置资源提供者
func (m *SkillManager) SetResourceProvider(skillName string, provider resources.ResourceProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.providers[skillName] = provider

	// 如果 Skill 已在缓存中，更新其 Provider
	if skill, ok := m.cache[skillName]; ok {
		skill.Provider = provider
	}
}

// GetStore 获取底层的 SkillStore
func (m *SkillManager) GetStore() store.SkillStore {
	return m.store
}
