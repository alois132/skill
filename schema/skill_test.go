package schema

import (
	"context"
	"testing"

	"github.com/alois132/skill/schema/resources"
)

// 测试脚本1：简单字符串处理
func testScript1(ctx context.Context, input string) (string, error) {
	return "processed: " + input, nil
}

// 测试脚本2：map 数据处理
func testScript2(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"status":  "success",
		"message": "test completed",
		"data":    input,
	}
	return result, nil
}

// 测试错误脚本
func testErrorScript(ctx context.Context, input string) (string, error) {
	return "", context.DeadlineExceeded
}

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

func TestSkill_AutoExecute(t *testing.T) {
	// 创建测试脚本
	type stringInput struct {
		Value string `json:"value"`
	}
	type stringOutput struct {
		Message string `json:"message"`
	}

	type mapInput map[string]interface{}
	type mapOutput map[string]interface{}

	tests := []struct {
		name         string
		skill        *Skill
		args         string
		expectCount  int
	}{
		{
			name: "single script",
			skill: &Skill{
				Metadata: &SkillMetadata{Name: "test"},
				Body:     "使用<script>test1</script>测试",
				Scripts: []resources.Script{
					resources.NewEasyScript("test1", testScript1),
				},
			},
			args:        `{"value":"hello"}`,
			expectCount: 1,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := tt.skill.AutoExecute(ctx, tt.args)
			if err != nil {
				t.Errorf("AutoExecute() error = %v", err)
				return
			}

			if len(results) != tt.expectCount {
				t.Errorf("Expected %d results, got %d", tt.expectCount, len(results))
				return
			}
		})
	}
}

func TestSkill_AutoExecute_NoScripts(t *testing.T) {
	skill := &Skill{
		Body: "没有脚本的普通文本",
	}

	ctx := context.Background()
	_, err := skill.AutoExecute(ctx, "{}")
	if err == nil {
		t.Error("Expected error for no scripts, got nil")
	}
}

func TestSkill_AutoExecute_ScriptNotFound(t *testing.T) {
	skill := &Skill{
		Metadata: &SkillMetadata{Name: "test"},
		Body:     "执行<script>nonexistent</script>脚本",
		Scripts:  []resources.Script{},
	}

	ctx := context.Background()
	results, err := skill.AutoExecute(ctx, "{}")
	if err != nil {
		t.Errorf("AutoExecute() error = %v", err)
		return
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
		return
	}

	if results[0].Error == nil {
		t.Error("Expected error for missing script, got nil")
	}
}

func TestSkill_Execute(t *testing.T) {
	skill := &Skill{
		Metadata: &SkillMetadata{Name: "test", Description: "Test skill"},
		Body:     "测试<script>test1</script>和<script>test2</script>",
		Scripts: []resources.Script{
			resources.NewEasyScript("test1", testScript1),
			resources.NewEasyScript("test2", testScript2),
		},
	}

	ctx := context.Background()
	result, err := skill.Execute(ctx, `{"data":"test"}`)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
		return
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}
}
