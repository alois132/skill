package schema

import (
	"context"
	"testing"

	"github.com/alois132/skill/schema/resources"
)

func TestSkill_ParseXMLTags(t *testing.T) {
	skill := &Skill{
		Metadata: &SkillMetadata{
			Name:        "test_skill",
			Description: "Test skill",
		},
		Body: `
第一步：使用<script>init</script>初始化
第二步：使用<script>config</script>配置
参考：<reference>usage_guide</reference>
模板：<asset>template.png</asset>
`,
	}

	err := skill.ParseXMLTags()
	if err != nil {
		t.Fatalf("ParseXMLTags() error = %v", err)
	}

	if !skill.parsed {
		t.Error("Expected skill.parsed to be true")
	}

	if len(skill.parsedTags) != 4 {
		t.Errorf("Expected 4 parsed tags, got %d", len(skill.parsedTags))
	}
}

func TestSkill_GetScriptNames(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected []string
	}{
		{
			name:     "single script",
			body:     "使用<script>init</script>初始化",
			expected: []string{"init"},
		},
		{
			name:     "multiple scripts",
			body:     "<script>init</script><script>config</script><script>deploy</script>",
			expected: []string{"init", "config", "deploy"},
		},
		{
			name:     "no scripts",
			body:     "参考：<reference>guide</reference>",
			expected: nil,
		},
		{
			name:     "empty body",
			body:     "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skill := &Skill{Body: tt.body}
			result := skill.GetScriptNames()
			if len(result) != len(tt.expected) {
				t.Errorf("GetScriptNames() = %v, want %v", result, tt.expected)
				return
			}
			for i, name := range result {
				if name != tt.expected[i] {
					t.Errorf("GetScriptNames()[%d] = %v, want %v", i, name, tt.expected[i])
				}
			}
		})
	}
}

func TestSkill_GetReferenceNames(t *testing.T) {
	skill := &Skill{
		Body: `
参考1：<reference>guide1</reference>
参考2：<reference>guide2</reference>
脚本：<script>init</script>
`,
	}

	refNames := skill.GetReferenceNames()
	expected := []string{"guide1", "guide2"}

	if len(refNames) != len(expected) {
		t.Errorf("Expected %d references, got %d", len(expected), len(refNames))
		return
	}

	for i, name := range refNames {
		if name != expected[i] {
			t.Errorf("refNames[%d] = %s, want %s", i, name, expected[i])
		}
	}
}

func TestSkill_GetAssetNames(t *testing.T) {
	skill := &Skill{
		Body: `
模板1：<asset>template1.png</asset>
模板2：<asset>template2.png</asset>
脚本：<script>init</script>
`,
	}

	assetNames := skill.GetAssetNames()
	expected := []string{"template1.png", "template2.png"}

	if len(assetNames) != len(expected) {
		t.Errorf("Expected %d assets, got %d", len(expected), len(assetNames))
		return
	}

	for i, name := range assetNames {
		if name != expected[i] {
			t.Errorf("assetNames[%d] = %s, want %s", i, name, expected[i])
		}
	}
}

func TestSkill_HasXMLTags(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected bool
	}{
		{
			name:     "has script tag",
			body:     "<script>init</script>",
			expected: true,
		},
		{
			name:     "has reference tag",
			body:     "<reference>guide</reference>",
			expected: true,
		},
		{
			name:     "has asset tag",
			body:     "<asset>logo.png</asset>",
			expected: true,
		},
		{
			name:     "no tags",
			body:     "plain text",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skill := &Skill{Body: tt.body}
			result := skill.HasXMLTags()
			if result != tt.expected {
				t.Errorf("HasXMLTags() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSkill_WithProvider(t *testing.T) {
	ctx := context.Background()

	// 创建 Provider
	provider := resources.NewInlineProvider()
	provider.AddScript(resources.NewEasyScript("provider_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"source": "provider"}, nil
	}))
	provider.AddReference(&resources.Reference{Name: "provider_ref", Body: "Provider reference content"})

	// 创建带有 Provider 的 Skill
	skill := &Skill{
		Metadata: &SkillMetadata{
			Name:        "test_skill",
			Description: "Test skill with provider",
		},
		Body: "Test body with <script>provider_script</script> and <reference>provider_ref</reference>",
		Scripts: []resources.Script{
			resources.NewEasyScript("inline_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{"source": "inline"}, nil
			}),
		},
		References: []*resources.Reference{
			{Name: "inline_ref", Body: "Inline reference content"},
		},
		Provider: provider,
	}

	// 测试从 Provider 获取脚本
	result, err := skill.UseScript(ctx, "provider_script", `{}`)
	if err != nil {
		t.Fatalf("Failed to use provider script: %v", err)
	}
	if result != `{"source":"provider"}` {
		t.Errorf("Expected '{\"source\":\"provider\"}', got '%s'", result)
	}

	// 测试从内联获取脚本（Provider 中不存在）
	result, err = skill.UseScript(ctx, "inline_script", `{}`)
	if err != nil {
		t.Fatalf("Failed to use inline script: %v", err)
	}
	if result != `{"source":"inline"}` {
		t.Errorf("Expected '{\"source\":\"inline\"}', got '%s'", result)
	}

	// 测试从 Provider 获取参考文档
	ref, err := skill.ReadReference("provider_ref")
	if err != nil {
		t.Fatalf("Failed to read provider reference: %v", err)
	}
	if ref != "Provider reference content" {
		t.Errorf("Expected 'Provider reference content', got '%s'", ref)
	}

	// 测试从内联获取参考文档（Provider 中不存在）
	ref, err = skill.ReadReference("inline_ref")
	if err != nil {
		t.Fatalf("Failed to read inline reference: %v", err)
	}
	if ref != "Inline reference content" {
		t.Errorf("Expected 'Inline reference content', got '%s'", ref)
	}
}

func TestSkill_ProviderPriority(t *testing.T) {
	ctx := context.Background()

	// 创建 Provider，包含同名脚本
	provider := resources.NewInlineProvider()
	provider.AddScript(resources.NewEasyScript("shared_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"source": "provider"}, nil
	}))

	// 创建带有 Provider 和内联同名脚本的 Skill
	skill := &Skill{
		Metadata: &SkillMetadata{Name: "priority_test"},
		Body:     "Test body",
		Scripts: []resources.Script{
			resources.NewEasyScript("shared_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{"source": "inline"}, nil
			}),
		},
		Provider: provider,
	}

	// Provider 应该优先
	result, err := skill.UseScript(ctx, "shared_script", `{}`)
	if err != nil {
		t.Fatalf("Failed to use script: %v", err)
	}
	if result != `{"source":"provider"}` {
		t.Errorf("Expected '{\"source\":\"provider\"}', got '%s'", result)
	}
}

func TestSkill_NoProvider(t *testing.T) {
	ctx := context.Background()

	// 创建没有 Provider 的 Skill
	skill := &Skill{
		Metadata: &SkillMetadata{Name: "no_provider"},
		Body:     "Test body",
		Scripts: []resources.Script{
			resources.NewEasyScript("test_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{"result": "ok"}, nil
			}),
		},
		References: []*resources.Reference{
			{Name: "test_ref", Body: "Reference content"},
		},
	}

	// 测试内联脚本
	result, err := skill.UseScript(ctx, "test_script", `{}`)
	if err != nil {
		t.Fatalf("Failed to use script: %v", err)
	}
	if result != `{"result":"ok"}` {
		t.Errorf("Expected '{\"result\":\"ok\"}', got '%s'", result)
	}

	// 测试内联参考文档
	ref, err := skill.ReadReference("test_ref")
	if err != nil {
		t.Fatalf("Failed to read reference: %v", err)
	}
	if ref != "Reference content" {
		t.Errorf("Expected 'Reference content', got '%s'", ref)
	}

	// 测试不存在的脚本
	_, err = skill.UseScript(ctx, "non_existent", `{}`)
	if err == nil {
		t.Error("Expected error for non-existent script")
	}

	// 测试不存在的参考文档
	_, err = skill.ReadReference("non_existent")
	if err == nil {
		t.Error("Expected error for non-existent reference")
	}
}

func TestSkill_GlanceAndInspect(t *testing.T) {
	skill := &Skill{
		Metadata: &SkillMetadata{
			Name:        "test_skill",
			Description: "Test description",
		},
		Body: "Test body content",
	}

	// 测试 Glance
	glance := skill.Glance()
	if glance == "" {
		t.Error("Expected non-empty glance")
	}
	if glance != `{"name":"test_skill","description":"Test description"}` {
		t.Errorf("Unexpected glance: %s", glance)
	}

	// 测试 Inspect
	inspect := skill.Inspect()
	if inspect != "Test body content" {
		t.Errorf("Expected 'Test body content', got '%s'", inspect)
	}
}
