package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alois132/skill/schema"
)

// FileStore 基于文件系统的 Skill 存储实现
// 每个 Skill 存储为一个 JSON 文件
type FileStore struct {
	mu       sync.RWMutex
	basePath string
	config   *StoreConfig
}

// FileStoreOption FileStore 特有的配置选项
type FileStoreOption func(*FileStore)

// NewFileStore 创建一个新的文件系统 Skill 存储
// basePath: 存储 Skill 文件的根目录
func NewFileStore(basePath string, opts ...StoreOption) (*FileStore, error) {
	config := &StoreConfig{}
	for _, opt := range opts {
		opt(config)
	}

	// 确保目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", basePath, err)
	}

	return &FileStore{
		basePath: basePath,
		config:   config,
	}, nil
}

// Get 从文件系统中获取指定名称的 Skill
func (s *FileStore) Get(ctx context.Context, name string) (*schema.Skill, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filePath := s.filePath(name)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("skill not found: " + name)
		}
		return nil, fmt.Errorf("failed to read skill file: %w", err)
	}

	var skill schema.Skill
	if err := json.Unmarshal(data, &skill); err != nil {
		return nil, fmt.Errorf("failed to unmarshal skill: %w", err)
	}

	return &skill, nil
}

// List 列出所有可用的 Skill 元数据
func (s *FileStore) List(ctx context.Context) ([]*schema.SkillMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	metadatas := make([]*schema.SkillMetadata, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(s.basePath, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // 跳过无法读取的文件
		}

		var skill schema.Skill
		if err := json.Unmarshal(data, &skill); err != nil {
			continue // 跳过无法解析的文件
		}

		if skill.Metadata != nil {
			metadatas = append(metadatas, skill.Metadata)
		}
	}

	return metadatas, nil
}

// Put 保存 Skill 到文件系统
func (s *FileStore) Put(ctx context.Context, skill *schema.Skill) error {
	if skill == nil {
		return errors.New("skill cannot be nil")
	}
	if skill.Metadata == nil || skill.Metadata.Name == "" {
		return errors.New("skill metadata name cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := s.filePath(skill.Metadata.Name)
	data, err := json.MarshalIndent(skill, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal skill: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write skill file: %w", err)
	}

	return nil
}

// Delete 从文件系统中删除指定名称的 Skill
func (s *FileStore) Delete(ctx context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := s.filePath(name)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New("skill not found: " + name)
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete skill file: %w", err)
	}

	return nil
}

// Exists 检查指定名称的 Skill 是否存在
func (s *FileStore) Exists(ctx context.Context, name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filePath := s.filePath(name)
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// filePath 生成 Skill 文件的完整路径
func (s *FileStore) filePath(name string) string {
	key := name
	if s.config.Namespace != "" {
		key = s.config.Namespace + "_" + name
	}
	return filepath.Join(s.basePath, key+".json")
}

// GetBasePath 获取存储的根目录路径
func (s *FileStore) GetBasePath() string {
	return s.basePath
}

// Ensure FileStore implements SkillStore
var _ SkillStore = (*FileStore)(nil)
