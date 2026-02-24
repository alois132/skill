package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/alois132/skill/schema"
)

func TestFileStore_Get(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file store: %v", err)
	}

	// 测试获取不存在的 Skill
	_, err = store.Get(ctx, "non_existent")
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

func TestFileStore_List(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file store: %v", err)
	}

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

func TestFileStore_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file store: %v", err)
	}

	// 删除不存在的 Skill
	err = store.Delete(ctx, "non_existent")
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

func TestFileStore_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file store: %v", err)
	}

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

func TestFileStore_WithNamespace(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	store, err := NewFileStore(tmpDir, WithNamespace("myapp"))
	if err != nil {
		t.Fatalf("Failed to create file store: %v", err)
	}

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

	// 验证文件使用了 namespace 前缀
	expectedFile := filepath.Join(tmpDir, "myapp_test.json")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist", expectedFile)
	}
}

func TestFileStore_InvalidPath(t *testing.T) {
	// 尝试在无效路径创建 store
	_, err := NewFileStore("/invalid/path/that/cannot/be/created")
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

func TestFileStore_InvalidSkill(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file store: %v", err)
	}

	// 测试 nil Skill
	err = store.Put(ctx, nil)
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
}

func TestFileStore_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	// 创建 store 并添加 skill
	store1, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file store: %v", err)
	}

	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{
			Name:        "persistent_skill",
			Description: "This skill should persist",
		},
		Body: "Persistent body",
	}
	if err := store1.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	// 创建新的 store 实例（模拟重启）
	store2, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create second file store: %v", err)
	}

	// 应该能够读取之前保存的 skill
	loaded, err := store2.Get(ctx, "persistent_skill")
	if err != nil {
		t.Fatalf("Failed to get skill from new store: %v", err)
	}
	if loaded.Metadata.Description != "This skill should persist" {
		t.Errorf("Expected description 'This skill should persist', got '%s'", loaded.Metadata.Description)
	}
}
