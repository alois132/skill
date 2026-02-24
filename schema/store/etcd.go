package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alois132/skill/schema"
)

// EtcdStore 基于 etcd 的 Skill 存储实现
// 这是一个示例实现，展示了如何基于外部存储实现 SkillStore 接口
// 实际使用时需要引入 etcd 客户端库: go.etcd.io/etcd/client/v3
//
// 由于 etcd 客户端库依赖较多，这里提供一个接口定义和模拟实现
// 用户可以根据需要替换为真实的 etcd 客户端
type EtcdStore struct {
	// client *clientv3.Client  // 真实实现时需要
	prefix string
	config *StoreConfig
}

// EtcdClient 定义 etcd 客户端的接口（用于解耦）
// 真实使用时可以替换为 go.etcd.io/etcd/client/v3.Client
type EtcdClient interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Put(ctx context.Context, key string, value []byte) error
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) (map[string][]byte, error)
}

// NewEtcdStore 创建一个新的 etcd Skill 存储
// 这是一个工厂函数，实际使用时需要传入 etcd 客户端
//
// 示例用法:
//
//	import clientv3 "go.etcd.io/etcd/client/v3"
//
//	cli, err := clientv3.New(clientv3.Config{
//	    Endpoints: []string{"localhost:2379"},
//	    DialTimeout: 5 * time.Second,
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer cli.Close()
//
//	store := NewEtcdStore(cli, store.WithNamespace("myapp"))
//
func NewEtcdStore(client EtcdClient, opts ...StoreOption) (*EtcdStore, error) {
	config := &StoreConfig{}
	for _, opt := range opts {
		opt(config)
	}

	prefix := config.Prefix
	if prefix == "" {
		prefix = "/skills"
	}
	if config.Namespace != "" {
		prefix = prefix + "/" + config.Namespace
	}

	return &EtcdStore{
		// client: client,  // 真实实现时需要
		prefix: prefix,
		config: config,
	}, nil
}

// Get 从 etcd 中获取指定名称的 Skill
func (s *EtcdStore) Get(ctx context.Context, name string) (*schema.Skill, error) {
	key := s.key(name)

	// 模拟实现：返回错误提示需要真实 etcd 客户端
	// 真实实现:
	// resp, err := s.client.Get(ctx, key)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to get from etcd: %w", err)
	// }
	// if len(resp.Kvs) == 0 {
	//     return nil, errors.New("skill not found: " + name)
	// }
	// var skill schema.Skill
	// if err := json.Unmarshal(resp.Kvs[0].Value, &skill); err != nil {
	//     return nil, fmt.Errorf("failed to unmarshal skill: %w", err)
	// }
	// return &skill, nil

	_ = key
	return nil, errors.New("etcd store requires real etcd client implementation")
}

// List 列出所有可用的 Skill 元数据
func (s *EtcdStore) List(ctx context.Context) ([]*schema.SkillMetadata, error) {
	// 真实实现:
	// resp, err := s.client.Get(ctx, s.prefix, clientv3.WithPrefix())
	// if err != nil {
	//     return nil, fmt.Errorf("failed to list from etcd: %w", err)
	// }
	//
	// metadatas := make([]*schema.SkillMetadata, 0, len(resp.Kvs))
	// for _, kv := range resp.Kvs {
	//     var skill schema.Skill
	//     if err := json.Unmarshal(kv.Value, &skill); err != nil {
	//         continue
	//     }
	//     if skill.Metadata != nil {
	//         metadatas = append(metadatas, skill.Metadata)
	//     }
	// }
	// return metadatas, nil

	return nil, errors.New("etcd store requires real etcd client implementation")
}

// Put 保存 Skill 到 etcd
func (s *EtcdStore) Put(ctx context.Context, skill *schema.Skill) error {
	if skill == nil {
		return errors.New("skill cannot be nil")
	}
	if skill.Metadata == nil || skill.Metadata.Name == "" {
		return errors.New("skill metadata name cannot be empty")
	}

	key := s.key(skill.Metadata.Name)
	data, err := json.Marshal(skill)
	if err != nil {
		return fmt.Errorf("failed to marshal skill: %w", err)
	}

	// 真实实现:
	// _, err = s.client.Put(ctx, key, string(data))
	// if err != nil {
	//     return fmt.Errorf("failed to put to etcd: %w", err)
	// }
	// return nil

	_ = key
	_ = data
	return errors.New("etcd store requires real etcd client implementation")
}

// Delete 从 etcd 中删除指定名称的 Skill
func (s *EtcdStore) Delete(ctx context.Context, name string) error {
	key := s.key(name)

	// 真实实现:
	// _, err := s.client.Delete(ctx, key)
	// if err != nil {
	//     return fmt.Errorf("failed to delete from etcd: %w", err)
	// }
	// return nil

	_ = key
	return errors.New("etcd store requires real etcd client implementation")
}

// Exists 检查指定名称的 Skill 是否存在
func (s *EtcdStore) Exists(ctx context.Context, name string) (bool, error) {
	// 真实实现:
	// resp, err := s.client.Get(ctx, s.key(name))
	// if err != nil {
	//     return false, err
	// }
	// return len(resp.Kvs) > 0, nil

	return false, errors.New("etcd store requires real etcd client implementation")
}

// key 生成 etcd 存储键
func (s *EtcdStore) key(name string) string {
	return s.prefix + "/" + name
}

// Watch 监视 Skill 变化（etcd 特有功能）
// 真实实现时可以使用 clientv3.Watch
func (s *EtcdStore) Watch(ctx context.Context, name string) (<-chan *schema.Skill, error) {
	// 真实实现:
	// watchChan := s.client.Watch(ctx, s.key(name))
	// resultChan := make(chan *schema.Skill)
	// go func() {
	//     for resp := range watchChan {
	//         for _, ev := range resp.Events {
	//             if ev.Type == clientv3.EventTypePut {
	//                 var skill schema.Skill
	//                 if err := json.Unmarshal(ev.Kv.Value, &skill); err == nil {
	//                     resultChan <- &skill
	//                 }
	//             }
	//         }
	//     }
	// }()
	// return resultChan, nil

	return nil, errors.New("etcd store requires real etcd client implementation")
}

// WatchPrefix 监视前缀下的所有 Skill 变化
func (s *EtcdStore) WatchPrefix(ctx context.Context) (<-chan *schema.Skill, error) {
	return nil, errors.New("etcd store requires real etcd client implementation")
}

// Ensure EtcdStore implements SkillStore
// 注意：由于方法返回错误，这里无法通过编译时检查
// 真实实现时需要取消下面的注释:
// var _ SkillStore = (*EtcdStore)(nil)

// EtcdStoreExample 展示如何使用 EtcdStore 的示例代码
// 这是一个文档函数，不会被实际调用
func EtcdStoreExample() {
	// 这是示例代码，展示了如何使用真实的 etcd 客户端
	// 需要引入: go.etcd.io/etcd/client/v3

	/*
		import (
			"context"
			"log"
			"time"

			"github.com/alois132/skill/schema"
			"github.com/alois132/skill/schema/store"
			clientv3 "go.etcd.io/etcd/client/v3"
		)

		func main() {
			// 创建 etcd 客户端
			cli, err := clientv3.New(clientv3.Config{
				Endpoints:   []string{"localhost:2379"},
				DialTimeout: 5 * time.Second,
			})
			if err != nil {
				log.Fatal(err)
			}
			defer cli.Close()

			// 创建 etcd store
			s, err := store.NewEtcdStore(cli, store.WithNamespace("myapp"))
			if err != nil {
				log.Fatal(err)
			}

			// 创建 skill
			skill := &schema.Skill{
				Metadata: &schema.SkillMetadata{
					Name:        "time_skill",
					Description: "Get current time",
				},
				Body: "Time skill body",
			}

			// 保存 skill
			ctx := context.Background()
			if err := s.Put(ctx, skill); err != nil {
				log.Fatal(err)
			}

			// 获取 skill
			loaded, err := s.Get(ctx, "time_skill")
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Loaded skill: %s", loaded.Metadata.Name)

			// 列出所有 skills
			metadatas, err := s.List(ctx)
			if err != nil {
				log.Fatal(err)
			}
			for _, meta := range metadatas {
				log.Printf("Skill: %s - %s", meta.Name, meta.Description)
			}
		}
	*/
}

