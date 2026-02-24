package store

import (
	"context"
	"errors"
	"sync"

	"github.com/alois132/skill/schema"
	"github.com/alois132/skill/schema/resources"
)

// MemoryStore 内存中的 Skill 存储实现
// 适用于测试和开发环境，数据不会持久化
type MemoryStore struct {
	mu     sync.RWMutex
	skills map[string]*schema.Skill
	config *StoreConfig
}

// NewMemoryStore 创建一个新的内存 Skill 存储
func NewMemoryStore(opts ...StoreOption) *MemoryStore {
	config := &StoreConfig{}
	for _, opt := range opts {
		opt(config)
	}

	return &MemoryStore{
		skills: make(map[string]*schema.Skill),
		config: config,
	}
}

// Get 从内存中获取指定名称的 Skill
func (s *MemoryStore) Get(ctx context.Context, name string) (*schema.Skill, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.key(name)
	skill, ok := s.skills[key]
	if !ok {
		return nil, errors.New("skill not found: " + name)
	}

	// 返回副本以避免外部修改
	return s.copySkill(skill), nil
}

// List 列出所有可用的 Skill 元数据
func (s *MemoryStore) List(ctx context.Context) ([]*schema.SkillMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metadatas := make([]*schema.SkillMetadata, 0, len(s.skills))
	for _, skill := range s.skills {
		if skill.Metadata != nil {
			metadatas = append(metadatas, skill.Metadata)
		}
	}
	return metadatas, nil
}

// Put 保存 Skill 到内存
func (s *MemoryStore) Put(ctx context.Context, skill *schema.Skill) error {
	if skill == nil {
		return errors.New("skill cannot be nil")
	}
	if skill.Metadata == nil || skill.Metadata.Name == "" {
		return errors.New("skill metadata name cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.key(skill.Metadata.Name)
	s.skills[key] = s.copySkill(skill)
	return nil
}

// Delete 从内存中删除指定名称的 Skill
func (s *MemoryStore) Delete(ctx context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.key(name)
	if _, ok := s.skills[key]; !ok {
		return errors.New("skill not found: " + name)
	}

	delete(s.skills, key)
	return nil
}

// Exists 检查指定名称的 Skill 是否存在
func (s *MemoryStore) Exists(ctx context.Context, name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.key(name)
	_, ok := s.skills[key]
	return ok, nil
}

// key 生成存储键
func (s *MemoryStore) key(name string) string {
	if s.config.Namespace != "" {
		return s.config.Namespace + "/" + name
	}
	return name
}

// copySkill 创建 Skill 的浅拷贝
// 注意：由于 Scripts 是接口类型，这里只做浅拷贝
// 适用于只读场景，如果需要完全隔离，需要更复杂的深拷贝逻辑
func (s *MemoryStore) copySkill(skill *schema.Skill) *schema.Skill {
	if skill == nil {
		return nil
	}

	// 创建新的 Skill 实例
	copied := &schema.Skill{
		Metadata: skill.Metadata,
		Body:     skill.Body,
	}

	// 拷贝 Scripts 切片（浅拷贝，元素是接口）
	if skill.Scripts != nil {
		copied.Scripts = make([]resources.Script, len(skill.Scripts))
		for i, script := range skill.Scripts {
			copied.Scripts[i] = script
		}
	}

	// 拷贝 References 切片
	if skill.References != nil {
		copied.References = make([]*resources.Reference, len(skill.References))
		for i, ref := range skill.References {
			copied.References[i] = ref
		}
	}

	// 拷贝 Assets 切片
	if skill.Assets != nil {
		copied.Assets = make([]*resources.Asset, len(skill.Assets))
		for i, asset := range skill.Assets {
			copied.Assets[i] = asset
		}
	}

	return copied
}

// GetAll 获取所有 Skills（仅用于测试）
func (s *MemoryStore) GetAll() map[string]*schema.Skill {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*schema.Skill, len(s.skills))
	for k, v := range s.skills {
		result[k] = s.copySkill(v)
	}
	return result
}

// Clear 清空所有 Skills（仅用于测试）
func (s *MemoryStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.skills = make(map[string]*schema.Skill)
}

// Ensure MemoryStore implements SkillStore
var _ SkillStore = (*MemoryStore)(nil)
