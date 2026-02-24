package store

import (
	"context"
	"github.com/alois132/skill/schema"
)

// SkillStore 定义 Skill 的持久化存储接口
// 实现此接口可以将 Skill 存储到各种后端（内存、文件、etcd、数据库等）
type SkillStore interface {
	// Get 从存储中获取指定名称的 Skill
	// 如果 Skill 不存在，返回 error
	Get(ctx context.Context, name string) (*schema.Skill, error)

	// List 列出所有可用的 Skill 元数据
	// 用于快速浏览可用的 Skills 而不加载完整内容
	List(ctx context.Context) ([]*schema.SkillMetadata, error)

	// Put 保存 Skill 到存储
	// 如果 Skill 已存在，则更新；否则创建
	Put(ctx context.Context, skill *schema.Skill) error

	// Delete 从存储中删除指定名称的 Skill
	// 如果 Skill 不存在，返回 error
	Delete(ctx context.Context, name string) error

	// Exists 检查指定名称的 Skill 是否存在
	Exists(ctx context.Context, name string) (bool, error)
}

// StoreOption SkillStore 的配置选项
type StoreOption func(*StoreConfig)

// StoreConfig SkillStore 的配置
type StoreConfig struct {
	Namespace string // 命名空间，用于隔离不同环境的 Skill
	Prefix    string // 键前缀
}

// WithNamespace 设置命名空间
func WithNamespace(ns string) StoreOption {
	return func(c *StoreConfig) {
		c.Namespace = ns
	}
}

// WithPrefix 设置键前缀
func WithPrefix(prefix string) StoreOption {
	return func(c *StoreConfig) {
		c.Prefix = prefix
	}
}
