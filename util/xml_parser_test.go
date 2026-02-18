package util

import (
	"reflect"
	"testing"
)

func TestParseXMLTags(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected []XMLTag
	}{
		{
			name: "single script tag",
			body: "第一步：使用<script>init_skill</script>初始化skill",
			expected: []XMLTag{
				{TagName: "script", Content: "init_skill"},
			},
		},
		{
			name: "multiple script tags",
			body: `第一步：<script>init</script>
第二步：<script>config</script>`,
			expected: []XMLTag{
				{TagName: "script", Content: "init"},
				{TagName: "script", Content: "config"},
			},
		},
		{
			name: "mixed tags",
			body: `参考：<reference>usage_guide</reference>
脚本：<script>init</script>
模板：<asset>template.png</asset>`,
			expected: []XMLTag{
				{TagName: "reference", Content: "usage_guide"},
				{TagName: "script", Content: "init"},
				{TagName: "asset", Content: "template.png"},
			},
		},
		{
			name:     "empty body",
			body:     "",
			expected: nil,
		},
		{
			name:     "no tags",
			body:     "这是一段普通文本，没有 XML 标记",
			expected: nil,
		},
		{
			name: "tags with whitespace",
			body: `<script>  init_skill  </script>`,
			expected: []XMLTag{
				{TagName: "script", Content: "init_skill"},
			},
		},
		{
			name: "multiple lines",
			body: `第一步：<script>
init_skill
</script>

参考文档：<reference>usage.md</reference>`,
			expected: []XMLTag{
				{TagName: "script", Content: "init_skill"},
				{TagName: "reference", Content: "usage.md"},
			},
		},
		{
			name: "complex content",
			body: `使用<script>complex_script</script>处理数据`,
			expected: []XMLTag{
				{TagName: "script", Content: "complex_script"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseXMLTags(tt.body)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseXMLTags() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractScriptNames(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected []string
	}{
		{
			name:     "single script",
			body:     "<script>init</script>",
			expected: []string{"init"},
		},
		{
			name:     "multiple scripts",
			body:     "<script>init</script><script>config</script><script>deploy</script>",
			expected: []string{"init", "config", "deploy"},
		},
		{
			name: "mixed with references",
			body: `<reference>guide</reference>
<script>init</script>
<reference>api</reference>`,
			expected: []string{"init"},
		},
		{
			name:     "no scripts",
			body:     "<reference>guide</reference><asset>logo.png</asset>",
			expected: nil,
		},
		{
			name:     "empty body",
			body:     "",
			expected: nil,
		},
		{
			name:     "no tags",
			body:     "plain text without tags",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractScriptNames(tt.body)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractScriptNames() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractReferenceNames(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected []string
	}{
		{
			name:     "single reference",
			body:     "<reference>guide</reference>",
			expected: []string{"guide"},
		},
		{
			name:     "multiple references",
			body:     "<reference>guide</reference><reference>api</reference>",
			expected: []string{"guide", "api"},
		},
		{
			name: "mixed with scripts",
			body: `<script>init</script>
<reference>guide</reference>`,
			expected: []string{"guide"},
		},
		{
			name:     "no references",
			body:     "<script>init</script><asset>logo.png</asset>",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractReferenceNames(tt.body)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractReferenceNames() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractAssetNames(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected []string
	}{
		{
			name:     "single asset",
			body:     "<asset>template.png</asset>",
			expected: []string{"template.png"},
		},
		{
			name:     "multiple assets",
			body:     "<asset>logo.png</asset><asset>banner.png</asset>",
			expected: []string{"logo.png", "banner.png"},
		},
		{
			name:     "no assets",
			body:     "<script>init</script><reference>guide</reference>",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractAssetNames(tt.body)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractAssetNames() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasXMLTags(t *testing.T) {
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
		{
			name:     "empty body",
			body:     "",
			expected: false,
		},
		{
			name:     "invalid tag format",
			body:     "<script>init",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasXMLTags(tt.body)
			if result != tt.expected {
				t.Errorf("HasXMLTags() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestRealWorldExample 测试真实场景
func TestRealWorldExample(t *testing.T) {
	body := `
第一步：使用<script>init_skill</script>初始化skill

此脚本将：
- 创建基础目录结构
- 生成配置文件
- 准备开发环境

参考：<reference>usage_guide</reference>
模板：<asset>template.png</asset>
`

	tags := ParseXMLTags(body)
	if len(tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(tags))
	}

	scriptNames := ExtractScriptNames(body)
	if len(scriptNames) != 1 || scriptNames[0] != "init_skill" {
		t.Errorf("Expected [init_skill], got %v", scriptNames)
	}

	refNames := ExtractReferenceNames(body)
	if len(refNames) != 1 || refNames[0] != "usage_guide" {
		t.Errorf("Expected [usage_guide], got %v", refNames)
	}

	assetNames := ExtractAssetNames(body)
	if len(assetNames) != 1 || assetNames[0] != "template.png" {
		t.Errorf("Expected [template.png], got %v", assetNames)
	}
}
