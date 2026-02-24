package store

import (
	"context"
	"testing"

	"github.com/alois132/skill/schema"
	"github.com/alois132/skill/schema/resources"
)

func TestMemoryStore_Get(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	// 测试获取不存在的 Skill
	_, err := store.Get(ctx, "non_existent")
	if err == nil {
		t.Error("Expected error for non-existent skill")
	}

	// 添加一个 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{
			Name:        "test_skill",
			Description: "Test skill",
		},
		Body: "Test body",
	}
	if err := store.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	// 获取存在的 Skill
	loaded, err := store.Get(ctx, "test_skill")
	if err != nil {
		t.Fatalf("Failed to get skill: %v", err)
	}
	if loaded.Metadata.Name != "test_skill" {
		t.Errorf("Expected skill name 'test_skill', got '%s'", loaded.Metadata.Name)
	}
}

func TestMemoryStore_List(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	// 空列表
	metadatas, err := store.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list: %v", err)
	}
	if len(metadatas) != 0 {
		t.Errorf("Expected 0 skills, got %d", len(metadatas))
	}

	// 添加多个 Skills
	skills := []*schema.Skill{
		{
			Metadata: &schema.SkillMetadata{Name: "skill1", Description: "Skill 1"},
			Body:     "Body 1",
		},
		{
			Metadata: &schema.SkillMetadata{Name: "skill2", Description: "Skill 2"},
			Body:     "Body 2",
		},
	}

	for _, skill := range skills {
		if err := store.Put(ctx, skill); err != nil {
			t.Fatalf("Failed to put skill: %v", err)
		}
	}

	// 列出所有 Skills
	metadatas, err = store.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list: %v", err)
	}
	if len(metadatas) != 2 {
		t.Errorf("Expected 2 skills, got %d", len(metadatas))
	}
}

func TestMemoryStore_Put(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	// 测试 nil Skill
	err := store.Put(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil skill")
	}

	// 测试空名称 Skill
	err = store.Put(ctx, &schema.Skill{
		Metadata: &schema.SkillMetadata{Name: ""},
	})
	if err == nil {
		t.Error("Expected error for empty skill name")
	}

	// 测试更新已存在的 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{Name: "test", Description: "Original"},
		Body:     "Original body",
	}
	if err := store.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	updated := &schema.Skill{
		Metadata: &schema.SkillMetadata{Name: "test", Description: "Updated"},
		Body:     "Updated body",
	}
	if err := store.Put(ctx, updated); err != nil {
		t.Fatalf("Failed to update skill: %v", err)
	}

	loaded, err := store.Get(ctx, "test")
	if err != nil {
		t.Fatalf("Failed to get skill: %v", err)
	}
	if loaded.Metadata.Description != "Updated" {
		t.Errorf("Expected description 'Updated', got '%s'", loaded.Metadata.Description)
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	// 删除不存在的 Skill
	err := store.Delete(ctx, "non_existent")
	if err == nil {
		t.Error("Expected error for non-existent skill")
	}

	// 添加并删除 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{Name: "to_delete"},
		Body:     "Body",
	}
	if err := store.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	if err := store.Delete(ctx, "to_delete"); err != nil {
		t.Fatalf("Failed to delete skill: %v", err)
	}

	// 确认已删除
	_, err = store.Get(ctx, "to_delete")
	if err == nil {
		t.Error("Expected error after deletion")
	}
}

func TestMemoryStore_Exists(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	// 检查不存在的 Skill
	exists, err := store.Exists(ctx, "non_existent")
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if exists {
		t.Error("Expected skill to not exist")
	}

	// 添加 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{Name: "test"},
		Body:     "Body",
	}
	if err := store.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	// 检查存在的 Skill
	exists, err = store.Exists(ctx, "test")
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if !exists {
		t.Error("Expected skill to exist")
	}
}

func TestMemoryStore_WithNamespace(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore(WithNamespace("myapp"))

	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{Name: "test"},
		Body:     "Body",
	}
	if err := store.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	// 使用相同的 namespace 应该能找到
	loaded, err := store.Get(ctx, "test")
	if err != nil {
		t.Fatalf("Failed to get skill: %v", err)
	}
	if loaded.Metadata.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", loaded.Metadata.Name)
	}

	// 使用不同的 namespace 应该找不到
	store2 := NewMemoryStore(WithNamespace("other"))
	_, err = store2.Get(ctx, "test")
	if err == nil {
		t.Error("Expected error for different namespace")
	}
}

func TestMemoryStore_Clear(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	// 添加 Skills
	for i := 0; i < 3; i++ {
		skill := &schema.Skill{
			Metadata: &schema.SkillMetadata{Name: "skill" + string(rune('0'+i))},
			Body:     "Body",
		}
		if err := store.Put(ctx, skill); err != nil {
			t.Fatalf("Failed to put skill: %v", err)
		}
	}

	// 清空
	store.Clear()

	// 确认已清空
	all := store.GetAll()
	if len(all) != 0 {
		t.Errorf("Expected 0 skills after clear, got %d", len(all))
	}
}

func TestMemoryStore_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	// 并发写入
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			skill := &schema.Skill{
				Metadata: &schema.SkillMetadata{Name: "skill" + string(rune('0'+n))},
				Body:     "Body",
			}
			store.Put(ctx, skill)
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证所有 Skills 都已写入
	all := store.GetAll()
	if len(all) != 10 {
		t.Errorf("Expected 10 skills, got %d", len(all))
	}
}

func TestMemoryStore_WithScriptsAndReferences(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	// 创建带有 Scripts 和 References 的 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{
			Name:        "test_skill",
			Description: "Test with resources",
		},
		Body: "Test body with <script>test_script</script>",
		Scripts: []resources.Script{
			resources.NewEasyScript("test_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{"result": "ok"}, nil
			}),
		},
		References: []*resources.Reference{
			{Name: "ref1", Body: "Reference content"},
		},
	}

	if err := store.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	loaded, err := store.Get(ctx, "test_skill")
	if err != nil {
		t.Fatalf("Failed to get skill: %v", err)
	}

	if len(loaded.Scripts) != 1 {
		t.Errorf("Expected 1 script, got %d", len(loaded.Scripts))
	}

	if len(loaded.References) != 1 {
		t.Errorf("Expected 1 reference, got %d", len(loaded.References))
	}
}
