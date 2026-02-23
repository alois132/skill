package eino

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/alois132/skill/core"
	"github.com/alois132/skill/schema"
)

// 创建测试用的 time_skill
func createTestTimeSkill() *schema.Skill {
	type TimeInput struct {
		Format   string `json:"format"`
		Timezone string `json:"timezone"`
	}
	type TimeOutput struct {
		Time string `json:"time"`
		Unix int64  `json:"unix"`
	}

	getCurrentTime := func(ctx context.Context, input TimeInput) (TimeOutput, error) {
		return TimeOutput{
			Time: "2024-01-01 12:00:00",
			Unix: 1704100800,
		}, nil
	}

	getTimezone := func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"timezone": "UTC",
		}, nil
	}

	return core.CreateSkill(
		"time_skill",
		"Get current time in various formats and timezone information",
		core.WithScript(core.CreateScript("get_current_time", getCurrentTime)),
		core.WithScript(core.CreateScript("get_timezone", getTimezone)),
		core.WithAutoParsedBody(`
获取当前时间的 Skill

## 使用方法

### 1. 获取当前时间
使用 <script>get_current_time</script> 脚本获取当前时间。

输入参数：
- format: 时间格式
  - "iso" - ISO 8601 格式 (默认)
  - "local" - 本地时间格式

### 2. 获取时区信息
使用 <script>get_timezone</script> 脚本获取时区信息。

## 参考文档
更多时间格式说明请参考：<reference>time_format_guide</reference>
`),
		core.WithReference("time_format_guide", `# 时间格式指南

## Go 时间布局参考
Go 使用以下参考时间进行格式化：
  "Mon Jan 2 15:04:05 MST 2006"

常用格式：
- "2006-01-02" - 日期
- "15:04:05" - 时间
`),
	)
}

func TestSkillTool(t *testing.T) {
	ctx := context.Background()
	skill := createTestTimeSkill()
	tool := NewSkillTool(skill)

	// 测试 Info
	t.Run("Info", func(t *testing.T) {
		info, err := tool.Info(ctx)
		if err != nil {
			t.Fatalf("Info() error = %v", err)
		}
		if info.Name != "time_skill" {
			t.Errorf("Info().Name = %v, want %v", info.Name, "time_skill")
		}
		if info.Desc != "Get current time in various formats and timezone information" {
			t.Errorf("Info().Desc = %v, want %v", info.Desc, "Get current time in various formats and timezone information")
		}
		// SkillTool 不需要参数
		if info.ParamsOneOf != nil {
			t.Error("Info().ParamsOneOf should be nil for SkillTool")
		}
	})

	// 测试 InvokableRun
	t.Run("InvokableRun", func(t *testing.T) {
		result, err := tool.InvokableRun(ctx, "")
		if err != nil {
			t.Fatalf("InvokableRun() error = %v", err)
		}
		if result == "" {
			t.Error("InvokableRun() returned empty body")
		}
		// 验证返回的内容包含预期的标记
		if !contains(result, "<script>get_current_time</script>") {
			t.Error("InvokableRun() result should contain <script>get_current_time</script>")
		}
		if !contains(result, "<reference>time_format_guide</reference>") {
			t.Error("InvokableRun() result should contain <reference>time_format_guide</reference>")
		}
	})
}

func TestUseScriptTool(t *testing.T) {
	ctx := context.Background()
	skill := createTestTimeSkill()
	tool := NewUseScriptTool(skill)

	// 测试 Info
	t.Run("Info", func(t *testing.T) {
		info, err := tool.Info(ctx)
		if err != nil {
			t.Fatalf("Info() error = %v", err)
		}
		if info.Name != "use_script" {
			t.Errorf("Info().Name = %v, want %v", info.Name, "use_script")
		}
		if info.ParamsOneOf == nil {
			t.Error("Info().ParamsOneOf should not be nil for UseScriptTool")
		}
	})

	// 测试 InvokableRun - 成功执行脚本
	t.Run("InvokableRun_Success", func(t *testing.T) {
		args := UseScriptRequest{
			SkillName:  "time_skill",
			ScriptName: "get_current_time",
			Args:       `{"format":"local"}`,
		}
		argsJSON, _ := json.Marshal(args)

		result, err := tool.InvokableRun(ctx, string(argsJSON))
		if err != nil {
			t.Fatalf("InvokableRun() error = %v", err)
		}

		var output struct {
			Time string `json:"time"`
			Unix int64  `json:"unix"`
		}
		if err := json.Unmarshal([]byte(result), &output); err != nil {
			t.Fatalf("Failed to unmarshal result: %v", err)
		}
		if output.Time != "2024-01-01 12:00:00" {
			t.Errorf("Expected time 2024-01-01 12:00:00, got %s", output.Time)
		}
	})

	// 测试 InvokableRun - skill 不存在
	t.Run("InvokableRun_SkillNotFound", func(t *testing.T) {
		args := UseScriptRequest{
			SkillName:  "non_existent_skill",
			ScriptName: "get_current_time",
			Args:       `{}`,
		}
		argsJSON, _ := json.Marshal(args)

		_, err := tool.InvokableRun(ctx, string(argsJSON))
		if err == nil {
			t.Error("InvokableRun() should return error for non-existent skill")
		}
	})

	// 测试 InvokableRun - 无效的 JSON
	t.Run("InvokableRun_InvalidJSON", func(t *testing.T) {
		_, err := tool.InvokableRun(ctx, "invalid json")
		if err == nil {
			t.Error("InvokableRun() should return error for invalid JSON")
		}
	})
}

func TestReadReferenceTool(t *testing.T) {
	ctx := context.Background()
	skill := createTestTimeSkill()
	tool := NewReadReferenceTool(skill)

	// 测试 Info
	t.Run("Info", func(t *testing.T) {
		info, err := tool.Info(ctx)
		if err != nil {
			t.Fatalf("Info() error = %v", err)
		}
		if info.Name != "read_reference" {
			t.Errorf("Info().Name = %v, want %v", info.Name, "read_reference")
		}
		if info.ParamsOneOf == nil {
			t.Error("Info().ParamsOneOf should not be nil for ReadReferenceTool")
		}
	})

	// 测试 InvokableRun - 成功读取 reference
	t.Run("InvokableRun_Success", func(t *testing.T) {
		args := ReadReferenceRequest{
			SkillName:     "time_skill",
			ReferenceName: "time_format_guide",
		}
		argsJSON, _ := json.Marshal(args)

		result, err := tool.InvokableRun(ctx, string(argsJSON))
		if err != nil {
			t.Fatalf("InvokableRun() error = %v", err)
		}
		if !contains(result, "时间格式指南") {
			t.Error("InvokableRun() result should contain '时间格式指南'")
		}
	})

	// 测试 InvokableRun - skill 不存在
	t.Run("InvokableRun_SkillNotFound", func(t *testing.T) {
		args := ReadReferenceRequest{
			SkillName:     "non_existent_skill",
			ReferenceName: "time_format_guide",
		}
		argsJSON, _ := json.Marshal(args)

		_, err := tool.InvokableRun(ctx, string(argsJSON))
		if err == nil {
			t.Error("InvokableRun() should return error for non-existent skill")
		}
	})

	// 测试 InvokableRun - reference 不存在
	t.Run("InvokableRun_ReferenceNotFound", func(t *testing.T) {
		args := ReadReferenceRequest{
			SkillName:     "time_skill",
			ReferenceName: "non_existent_reference",
		}
		argsJSON, _ := json.Marshal(args)

		_, err := tool.InvokableRun(ctx, string(argsJSON))
		if err == nil {
			t.Error("InvokableRun() should return error for non-existent reference")
		}
	})
}

func TestToTools(t *testing.T) {
	skill := createTestTimeSkill()

	// 测试单个 skill
	t.Run("SingleSkill", func(t *testing.T) {
		tools := ToTools(skill)
		if len(tools) != 3 {
			t.Errorf("ToTools() returned %d tools, want 3", len(tools))
		}
	})

	// 测试多个 skills
	t.Run("MultipleSkills", func(t *testing.T) {
		skill2 := core.CreateSkill(
			"test_skill",
			"Test skill",
			core.WithBody("Test body"),
		)
		tools := ToTools(skill, skill2)
		if len(tools) != 4 {
			t.Errorf("ToTools() returned %d tools, want 4", len(tools))
		}
	})

	// 测试空 skills
	t.Run("EmptySkills", func(t *testing.T) {
		tools := ToTools()
		if tools != nil {
			t.Error("ToTools() should return nil for empty skills")
		}
	})
}

func TestToInvokableTools(t *testing.T) {
	skill := createTestTimeSkill()

	t.Run("SingleSkill", func(t *testing.T) {
		tools := ToInvokableTools(skill)
		if len(tools) != 3 {
			t.Errorf("ToInvokableTools() returned %d tools, want 3", len(tools))
		}
	})

	t.Run("EmptySkills", func(t *testing.T) {
		tools := ToInvokableTools()
		if tools != nil {
			t.Error("ToInvokableTools() should return nil for empty skills")
		}
	})
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
	if start+len(substr) > len(s) {
		return false
	}
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
