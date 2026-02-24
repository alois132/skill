package core

import (
	"context"
	"testing"

	"github.com/alois132/skill/schema"
	"github.com/alois132/skill/schema/resources"
	"github.com/alois132/skill/schema/store"
)

func TestSkillManager_GetSkill(t *testing.T) {
	ctx := context.Background()
	memStore := store.NewMemoryStore()
	manager := NewSkillManager(memStore)

	// 测试获取不存在的 Skill
	_, err := manager.GetSkill(ctx, "non_existent")
	if err == nil {
		t.Error("Expected error for non-existent skill")
	}

	// 添加 Skill 到 store
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{
			Name:        "test_skill",
			Description: "Test skill",
		},
		Body: "Test body",
	}
	if err := memStore.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	// 从 manager 获取
	loaded, err := manager.GetSkill(ctx, "test_skill")
	if err != nil {
		t.Fatalf("Failed to get skill: %v", err)
	}
	if loaded.Metadata.Name != "test_skill" {
		t.Errorf("Expected name 'test_skill', got '%s'", loaded.Metadata.Name)
	}

	// 第二次获取应该从缓存
	cached, err := manager.GetSkill(ctx, "test_skill")
	if err != nil {
		t.Fatalf("Failed to get skill from cache: %v", err)
	}
	if cached != loaded {
		t.Error("Expected cached skill to be the same object")
	}
}

func TestSkillManager_RegisterSkill(t *testing.T) {
	manager := NewSkillManager(nil)

	// 注册 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{
			Name:        "registered_skill",
			Description: "Registered skill",
		},
		Body: "Body",
	}

	if err := manager.RegisterSkill(skill); err != nil {
		t.Fatalf("Failed to register skill: %v", err)
	}

	// 验证已注册
	names := manager.GetCachedSkillNames()
	if len(names) != 1 || names[0] != "registered_skill" {
		t.Errorf("Expected ['registered_skill'], got %v", names)
	}

	// 测试 nil Skill
	err := manager.RegisterSkill(nil)
	if err == nil {
		t.Error("Expected error for nil skill")
	}

	// 测试空名称 Skill
	err = manager.RegisterSkill(&schema.Skill{
		Metadata: &schema.SkillMetadata{Name: ""},
	})
	if err == nil {
		t.Error("Expected error for empty skill name")
	}
}

func TestSkillManager_SaveSkill(t *testing.T) {
	ctx := context.Background()
	memStore := store.NewMemoryStore()
	manager := NewSkillManager(memStore)

	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{
			Name:        "save_test",
			Description: "Save test",
		},
		Body: "Body",
	}

	// 保存 Skill
	if err := manager.SaveSkill(ctx, skill); err != nil {
		t.Fatalf("Failed to save skill: %v", err)
	}

	// 验证已保存到 store
	loaded, err := memStore.Get(ctx, "save_test")
	if err != nil {
		t.Fatalf("Failed to get skill from store: %v", err)
	}
	if loaded.Metadata.Description != "Save test" {
		t.Errorf("Expected description 'Save test', got '%s'", loaded.Metadata.Description)
	}

	// 验证已缓存
	names := manager.GetCachedSkillNames()
	if len(names) != 1 {
		t.Errorf("Expected 1 cached skill, got %d", len(names))
	}
}

func TestSkillManager_ReloadSkill(t *testing.T) {
	ctx := context.Background()
	memStore := store.NewMemoryStore()
	manager := NewSkillManager(memStore)

	// 添加初始 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{
			Name:        "reload_test",
			Description: "Original",
		},
		Body: "Original body",
	}
	if err := memStore.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	// 首次获取（会缓存）
	original, err := manager.GetSkill(ctx, "reload_test")
	if err != nil {
		t.Fatalf("Failed to get skill: %v", err)
	}

	// 更新 store 中的 Skill
	updated := &schema.Skill{
		Metadata: &schema.SkillMetadata{
			Name:        "reload_test",
			Description: "Updated",
		},
		Body: "Updated body",
	}
	if err := memStore.Put(ctx, updated); err != nil {
		t.Fatalf("Failed to put updated skill: %v", err)
	}

	// 重新加载
	reloaded, err := manager.ReloadSkill(ctx, "reload_test")
	if err != nil {
		t.Fatalf("Failed to reload skill: %v", err)
	}

	// 验证已更新
	if reloaded.Metadata.Description != "Updated" {
		t.Errorf("Expected description 'Updated', got '%s'", reloaded.Metadata.Description)
	}

	// 验证缓存已更新
	cached, _ := manager.GetSkill(ctx, "reload_test")
	if cached.Metadata.Description != "Updated" {
		t.Errorf("Expected cached description 'Updated', got '%s'", cached.Metadata.Description)
	}

	// 验证是新对象
	if original == reloaded {
		t.Error("Expected reloaded skill to be a different object")
	}
}

func TestSkillManager_DeleteSkill(t *testing.T) {
	ctx := context.Background()
	memStore := store.NewMemoryStore()
	manager := NewSkillManager(memStore)

	// 添加 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{Name: "to_delete"},
		Body:     "Body",
	}
	if err := memStore.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	// 获取以缓存
	manager.GetSkill(ctx, "to_delete")

	// 删除
	if err := manager.DeleteSkill(ctx, "to_delete"); err != nil {
		t.Fatalf("Failed to delete skill: %v", err)
	}

	// 验证已从 store 删除
	_, err := memStore.Get(ctx, "to_delete")
	if err == nil {
		t.Error("Expected skill to be deleted from store")
	}

	// 验证已从缓存删除
	names := manager.GetCachedSkillNames()
	if len(names) != 0 {
		t.Errorf("Expected 0 cached skills, got %d", len(names))
	}
}

func TestSkillManager_ListSkills(t *testing.T) {
	ctx := context.Background()
	memStore := store.NewMemoryStore()
	manager := NewSkillManager(memStore)

	// 添加 Skills
	for i := 0; i < 3; i++ {
		skill := &schema.Skill{
			Metadata: &schema.SkillMetadata{
				Name:        "skill" + string(rune('0'+i)),
				Description: "Skill " + string(rune('0'+i)),
			},
			Body: "Body",
		}
		if err := memStore.Put(ctx, skill); err != nil {
			t.Fatalf("Failed to put skill: %v", err)
		}
	}

	// 列出 Skills
	metadatas, err := manager.ListSkills(ctx)
	if err != nil {
		t.Fatalf("Failed to list skills: %v", err)
	}
	if len(metadatas) != 3 {
		t.Errorf("Expected 3 skills, got %d", len(metadatas))
	}
}

func TestSkillManager_ListSkills_NoStore(t *testing.T) {
	manager := NewSkillManager(nil)

	// 直接注册 Skills
	for i := 0; i < 3; i++ {
		skill := &schema.Skill{
			Metadata: &schema.SkillMetadata{
				Name: "skill" + string(rune('0'+i)),
			},
			Body: "Body",
		}
		manager.RegisterSkill(skill)
	}

	// 列出 Skills（应该从缓存）
	ctx := context.Background()
	metadatas, err := manager.ListSkills(ctx)
	if err != nil {
		t.Fatalf("Failed to list skills: %v", err)
	}
	if len(metadatas) != 3 {
		t.Errorf("Expected 3 skills, got %d", len(metadatas))
	}
}

func TestSkillManager_UseScript(t *testing.T) {
	ctx := context.Background()
	memStore := store.NewMemoryStore()
	manager := NewSkillManager(memStore)

	// 创建带有脚本的 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{Name: "script_test"},
		Body:     "Test body",
		Scripts: []resources.Script{
			resources.NewEasyScript("test_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{"result": "success"}, nil
			}),
		},
	}
	if err := memStore.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	// 执行脚本
	result, err := manager.UseScript(ctx, "script_test", "test_script", `{}`)
	if err != nil {
		t.Fatalf("Failed to use script: %v", err)
	}
	if result != `{"result":"success"}` {
		t.Errorf("Expected '{\"result\":\"success\"}', got '%s'", result)
	}

	// 测试不存在的 Skill
	_, err = manager.UseScript(ctx, "non_existent", "test_script", `{}`)
	if err == nil {
		t.Error("Expected error for non-existent skill")
	}
}

func TestSkillManager_ReadReference(t *testing.T) {
	ctx := context.Background()
	memStore := store.NewMemoryStore()
	manager := NewSkillManager(memStore)

	// 创建带有参考文档的 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{Name: "ref_test"},
		Body:     "Test body",
		References: []*resources.Reference{
			{Name: "test_ref", Body: "Reference content"},
		},
	}
	if err := memStore.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	// 读取参考文档
	content, err := manager.ReadReference(ctx, "ref_test", "test_ref")
	if err != nil {
		t.Fatalf("Failed to read reference: %v", err)
	}
	if content != "Reference content" {
		t.Errorf("Expected 'Reference content', got '%s'", content)
	}

	// 测试不存在的 Skill
	_, err = manager.ReadReference(ctx, "non_existent", "test_ref")
	if err == nil {
		t.Error("Expected error for non-existent skill")
	}
}

func TestSkillManager_ClearCache(t *testing.T) {
	ctx := context.Background()
	memStore := store.NewMemoryStore()
	manager := NewSkillManager(memStore)

	// 添加并获取 Skills
	for i := 0; i < 3; i++ {
		skill := &schema.Skill{
			Metadata: &schema.SkillMetadata{Name: "skill" + string(rune('0'+i))},
			Body:     "Body",
		}
		if err := memStore.Put(ctx, skill); err != nil {
			t.Fatalf("Failed to put skill: %v", err)
		}
		manager.GetSkill(ctx, "skill"+string(rune('0'+i)))
	}

	// 验证缓存
	names := manager.GetCachedSkillNames()
	if len(names) != 3 {
		t.Errorf("Expected 3 cached skills, got %d", len(names))
	}

	// 清空缓存
	manager.ClearCache()

	// 验证已清空
	names = manager.GetCachedSkillNames()
	if len(names) != 0 {
		t.Errorf("Expected 0 cached skills after clear, got %d", len(names))
	}
}

func TestSkillManager_SetResourceProvider(t *testing.T) {
	ctx := context.Background()
	memStore := store.NewMemoryStore()
	manager := NewSkillManager(memStore)

	// 创建带有 Provider 的 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{Name: "provider_test"},
		Body:     "Test body",
	}
	if err := memStore.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	// 创建 Provider
	provider := resources.NewInlineProvider()
	provider.AddScript(resources.NewEasyScript("provider_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"from": "provider"}, nil
	}))

	// 设置 Provider
	manager.SetResourceProvider("provider_test", provider)

	// 获取 Skill 并验证 Provider 已设置
	loaded, err := manager.GetSkill(ctx, "provider_test")
	if err != nil {
		t.Fatalf("Failed to get skill: %v", err)
	}
	if loaded.Provider == nil {
		t.Error("Expected provider to be set")
	}

	// 执行 Provider 中的脚本
	result, err := manager.UseScript(ctx, "provider_test", "provider_script", `{}`)
	if err != nil {
		t.Fatalf("Failed to use script: %v", err)
	}
	if result != `{"from":"provider"}` {
		t.Errorf("Expected '{\"from\":\"provider\"}', got '%s'", result)
	}
}

func TestSkillManager_WithManagerResourceProvider(t *testing.T) {
	ctx := context.Background()
	memStore := store.NewMemoryStore()

	// 创建 Provider
	provider := resources.NewInlineProvider()
	provider.AddScript(resources.NewEasyScript("init_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"initialized": true}, nil
	}))

	// 创建带有 Provider 的 Manager
	manager := NewSkillManager(memStore, WithManagerResourceProvider("init_test", provider))

	// 添加 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{Name: "init_test"},
		Body:     "Test body",
	}
	if err := memStore.Put(ctx, skill); err != nil {
		t.Fatalf("Failed to put skill: %v", err)
	}

	// 获取 Skill，Provider 应该已设置
	loaded, err := manager.GetSkill(ctx, "init_test")
	if err != nil {
		t.Fatalf("Failed to get skill: %v", err)
	}
	if loaded.Provider == nil {
		t.Error("Expected provider to be set via option")
	}

	// 执行脚本
	result, err := manager.UseScript(ctx, "init_test", "init_script", `{}`)
	if err != nil {
		t.Fatalf("Failed to use script: %v", err)
	}
	if result != `{"initialized":true}` {
		t.Errorf("Expected '{\"initialized\":true}', got '%s'", result)
	}
}

func TestSkillManager_GetStore(t *testing.T) {
	memStore := store.NewMemoryStore()
	manager := NewSkillManager(memStore)

	if manager.GetStore() != memStore {
		t.Error("Expected GetStore to return the same store")
	}
}
