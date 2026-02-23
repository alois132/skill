package schema

import (
	"testing"
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
