package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/alois132/skill/core"
	"github.com/alois132/skill/schema/resources"
	"github.com/alois132/skill/schema/store"
)

// TestMemoryStoreExample 演示如何使用 MemoryStore
func TestMemoryStoreExample(t *testing.T) {
	ctx := context.Background()

	// 创建内存存储
	memStore := store.NewMemoryStore()

	// 创建 Skill
	skill := core.CreateSkill(
		"calculator",
		"A simple calculator skill",
		core.WithBody(`
Calculator Skill

Use <script>add</script> to add two numbers.
Use <script>subtract</script> to subtract two numbers.

See <reference>math_guide</reference> for more information.
`),
		core.WithScript(core.CreateScript("add", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
			a, _ := input["a"].(float64)
			b, _ := input["b"].(float64)
			return map[string]interface{}{"result": a + b}, nil
		})),
		core.WithScript(core.CreateScript("subtract", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
			a, _ := input["a"].(float64)
			b, _ := input["b"].(float64)
			return map[string]interface{}{"result": a - b}, nil
		})),
		core.WithReference("math_guide", "# Math Guide\n\nBasic arithmetic operations."),
	)

	// 保存 Skill 到存储
	if err := memStore.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	// 从存储获取 Skill
	loaded, err := memStore.Get(ctx, "calculator")
	if err != nil {
		t.Fatalf("Failed to get skill: %v", err)
	}

	fmt.Printf("Loaded skill: %s - %s\n", loaded.Metadata.Name, loaded.Metadata.Description)

	// 执行脚本
	result, err := loaded.UseScript(ctx, "add", `{"a":5,"b":3}`)
	if err != nil {
		t.Fatalf("Failed to use script: %v", err)
	}
	fmt.Printf("5 + 3 = %s\n", result)

	// 列出所有 Skills
	metadatas, err := memStore.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list skills: %v", err)
	}
	fmt.Printf("Total skills: %d\n", len(metadatas))
}

// TestFileStoreExample 演示如何使用 FileStore
// 注意：FileStore 只能存储可序列化的数据（Metadata 和 Body）
// Scripts 需要通过 ResourceProvider 动态加载
func TestFileStoreExample(t *testing.T) {
	ctx := context.Background()

	// 创建临时目录用于测试
	tmpDir := t.TempDir()

	// 创建文件存储
	fileStore, err := store.NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file store: %v", err)
	}

	// 创建 Skill（仅包含可序列化的 Metadata 和 Body）
	skill := core.CreateSkill(
		"greeting",
		"A greeting skill",
		core.WithBody(`
Greeting Skill

Use <script>say_hello</script> to greet someone.
`),
		// 注意：FileStore 不能存储本地函数脚本
		// 实际使用时，可以通过 ResourceProvider 动态加载脚本
	)

	// 保存 Skill（只保存 Metadata 和 Body）
	if err := fileStore.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	fmt.Println("Skill saved to file store")

	// 重新加载 Skill
	loaded, err := fileStore.Get(ctx, "greeting")
	if err != nil {
		t.Fatalf("Failed to get skill: %v", err)
	}

	fmt.Printf("Loaded skill: %s\n", loaded.Metadata.Name)
	fmt.Printf("Body: %s\n", loaded.Body)

	// 验证文件持久化
	exists, err := fileStore.Exists(ctx, "greeting")
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	fmt.Printf("Skill exists: %v\n", exists)
}

// TestSkillManagerExample 演示如何使用 SkillManager
func TestSkillManagerExample(t *testing.T) {
	ctx := context.Background()

	// 创建存储和管理器
	memStore := store.NewMemoryStore()
	manager := core.NewSkillManager(memStore)

	// 创建并注册 Skill
	skill := core.CreateSkill(
		"echo",
		"An echo skill",
		core.WithBody(`Echo Skill - use <script>echo</script>`),
		core.WithScript(core.CreateScript("echo", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
			return input, nil
		})),
	)

	// 保存到管理器
	if err := manager.SaveSkill(ctx, skill); err != nil {
		t.Fatalf("Failed to save skill: %v", err)
	}

	// 通过管理器执行脚本
	result, err := manager.UseScript(ctx, "echo", "echo", `{"test":"hello"}`)
	if err != nil {
		t.Fatalf("Failed to use script: %v", err)
	}
	fmt.Printf("Echo result: %s\n", result)

	// 列出所有 Skills
	metadatas, err := manager.ListSkills(ctx)
	if err != nil {
		t.Fatalf("Failed to list skills: %v", err)
	}
	fmt.Printf("Skills in manager: %d\n", len(metadatas))
	for _, meta := range metadatas {
		fmt.Printf("  - %s: %s\n", meta.Name, meta.Description)
	}
}

// TestNamespaceExample 演示如何使用命名空间隔离 Skills
func TestNamespaceExample(t *testing.T) {
	ctx := context.Background()

	// 创建两个不同命名空间的存储
	store1 := store.NewMemoryStore(store.WithNamespace("app1"))
	store2 := store.NewMemoryStore(store.WithNamespace("app2"))

	// 在 app1 中创建 Skill
	skill1 := core.CreateSkill(
		"shared",
		"App1 version",
		core.WithBody("App1 body"),
	)
	if err := store1.Put(ctx, skill1); err != nil {
		t.Fatalf("Failed to put skill to store1: %v", err)
	}

	// 在 app2 中创建同名 Skill
	skill2 := core.CreateSkill(
		"shared",
		"App2 version",
		core.WithBody("App2 body"),
	)
	if err := store2.Put(ctx, skill2); err != nil {
		t.Fatalf("Failed to put skill to store2: %v", err)
	}

	// 从各自的存储获取
	loaded1, _ := store1.Get(ctx, "shared")
	loaded2, _ := store2.Get(ctx, "shared")

	fmt.Printf("App1 skill: %s\n", loaded1.Metadata.Description)
	fmt.Printf("App2 skill: %s\n", loaded2.Metadata.Description)
}

// TestDynamicSkillUpdate 演示动态更新 Skill
func TestDynamicSkillUpdate(t *testing.T) {
	ctx := context.Background()

	memStore := store.NewMemoryStore()
	manager := core.NewSkillManager(memStore)

	// 创建初始版本
	skillV1 := core.CreateSkill(
		"versioned",
		"Version 1",
		core.WithBody("Initial version"),
	)
	if err := manager.SaveSkill(ctx, skillV1); err != nil {
		t.Fatalf("Failed to save skill v1: %v", err)
	}

	// 获取并验证
	loaded, _ := manager.GetSkill(ctx, "versioned")
	fmt.Printf("Initial: %s\n", loaded.Metadata.Description)

	// 更新 Skill
	skillV2 := core.CreateSkill(
		"versioned",
		"Version 2",
		core.WithBody("Updated version with new features"),
		core.WithScript(core.CreateScript("new_feature", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{"feature": "new"}, nil
		})),
	)
	if err := manager.SaveSkill(ctx, skillV2); err != nil {
		t.Fatalf("Failed to save skill v2: %v", err)
	}

	// 重新加载以获取最新版本
	reloaded, _ := manager.ReloadSkill(ctx, "versioned")
	fmt.Printf("Updated: %s\n", reloaded.Metadata.Description)
	fmt.Printf("Body: %s\n", reloaded.Body)
}

// TestResourceProviderWithStore 演示结合 ResourceProvider 和 Store 使用
func TestResourceProviderWithStore(t *testing.T) {
	ctx := context.Background()

	// 创建存储
	memStore := store.NewMemoryStore()

	// 创建资源提供者
	provider := core.CreateInlineProvider()
	provider.AddScript(core.CreateScript("dynamic_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"dynamic": true, "input": input}, nil
	}))
	provider.AddReference(&resources.Reference{
		Name: "dynamic_ref",
		Body: "This is dynamically loaded content",
	})

	// 创建使用 Provider 的 Skill
	skill := core.CreateSkill(
		"hybrid",
		"Hybrid skill with inline and dynamic resources",
		core.WithBody(`
Hybrid Skill

Inline script: <script>inline_script</script>
Dynamic script: <script>dynamic_script</script>
Dynamic ref: <reference>dynamic_ref</reference>
`),
		core.WithScript(core.CreateScript("inline_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{"source": "inline"}, nil
		})),
		core.WithResourceProvider(provider),
	)

	// 保存到存储
	if err := memStore.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	// 重新加载（注意：Provider 不会被序列化，需要重新设置）
	loaded, _ := memStore.Get(ctx, "hybrid")

	// 重新设置 Provider
	loaded.Provider = provider

	// 执行内联脚本
	result1, _ := loaded.UseScript(ctx, "inline_script", `{}`)
	fmt.Printf("Inline script: %s\n", result1)

	// 执行动态脚本（来自 Provider）
	result2, _ := loaded.UseScript(ctx, "dynamic_script", `{"key":"value"}`)
	fmt.Printf("Dynamic script: %s\n", result2)

	// 读取动态参考文档
	ref, _ := loaded.ReadReference("dynamic_ref")
	fmt.Printf("Dynamic reference: %s\n", ref)
}
