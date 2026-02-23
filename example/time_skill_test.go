package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// TestGetCurrentTime_ISOFormat 测试 ISO 格式时间获取
func TestGetCurrentTime_ISOFormat(t *testing.T) {
	ctx := context.Background()
	input := TimeInput{Format: "iso"}

	output, err := getCurrentTime(ctx, input)
	if err != nil {
		t.Fatalf("getCurrentTime() error = %v", err)
	}

	// 验证返回了时间字符串
	if output.Time == "" {
		t.Error("Expected non-empty time string")
	}

	// 验证 Unix 时间戳
	if output.Unix == 0 {
		t.Error("Expected non-zero unix timestamp")
	}

	// 验证 ISO 格式（应该包含 T 和时区信息）
	if !strings.Contains(output.Time, "T") {
		t.Errorf("ISO format should contain 'T', got: %s", output.Time)
	}

	t.Logf("ISO Time: %s, Unix: %d", output.Time, output.Unix)
}

// TestGetCurrentTime_LocalFormat 测试本地格式时间获取
func TestGetCurrentTime_LocalFormat(t *testing.T) {
	ctx := context.Background()
	input := TimeInput{Format: "local"}

	output, err := getCurrentTime(ctx, input)
	if err != nil {
		t.Fatalf("getCurrentTime() error = %v", err)
	}

	// 验证格式 "2006-01-02 15:04:05"
	expectedParts := 2
	actualParts := len(strings.Split(output.Time, " "))
	if actualParts != expectedParts {
		t.Errorf("Local format should have %d parts (date and time), got: %d", expectedParts, actualParts)
	}

	t.Logf("Local Time: %s", output.Time)
}

// TestGetCurrentTime_UnixFormat 测试 Unix 时间戳格式
func TestGetCurrentTime_UnixFormat(t *testing.T) {
	ctx := context.Background()
	input := TimeInput{Format: "unix"}

	before := time.Now().Unix()
	output, err := getCurrentTime(ctx, input)
	after := time.Now().Unix()

	if err != nil {
		t.Fatalf("getCurrentTime() error = %v", err)
	}

	// 验证 Unix 时间戳是数字字符串
	_ = output.Time // 确保使用变量
	if output.Time == "" {
		t.Errorf("Unix format should be a valid number string: %s", output.Time)
	}

	// 验证时间在合理范围内
	parsedTime := output.Unix
	if parsedTime < before || parsedTime > after {
		t.Errorf("Unix timestamp %d not in expected range [%d, %d]", parsedTime, before, after)
	}

	t.Logf("Unix Time: %s", output.Time)
}

// TestGetCurrentTime_CustomTimezone 测试自定义时区
func TestGetCurrentTime_CustomTimezone(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name     string
		timezone string
		wantErr  bool
	}{
		{
			name:     "UTC timezone",
			timezone: "UTC",
			wantErr:  false,
		},
		{
			name:     "Shanghai timezone",
			timezone: "Asia/Shanghai",
			wantErr:  false,
		},
		{
			name:     "Tokyo timezone",
			timezone: "Asia/Tokyo",
			wantErr:  false,
		},
		{
			name:     "New York timezone",
			timezone: "America/New_York",
			wantErr:  false,
		},
		{
			name:     "Invalid timezone",
			timezone: "Invalid/Timezone",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := TimeInput{
				Format:   "iso",
				Timezone: tc.timezone,
			}

			output, err := getCurrentTime(ctx, input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error for timezone %s, got nil", tc.timezone)
				}
				return
			}

			if err != nil {
				t.Fatalf("getCurrentTime() error = %v", err)
			}

			// 验证时区信息
			if output.Timezone != tc.timezone {
				t.Errorf("Expected timezone %s, got %s", tc.timezone, output.Timezone)
			}

			t.Logf("Timezone: %s, Time: %s", tc.timezone, output.Time)
		})
	}
}

// TestGetCurrentTime_CustomFormat 测试自定义格式
func TestGetCurrentTime_CustomFormat(t *testing.T) {
	ctx := context.Background()
	input := TimeInput{
		Format: "custom",
		Layout: "2006年01月02日 15:04",
	}

	output, err := getCurrentTime(ctx, input)
	if err != nil {
		t.Fatalf("getCurrentTime() error = %v", err)
	}

	// 验证包含中文"年"
	if !strings.Contains(output.Time, "年") {
		t.Errorf("Custom format should contain '年', got: %s", output.Time)
	}

	t.Logf("Custom Time: %s", output.Time)
}

// TestGetCurrentTime_DateOnly 测试仅日期格式
func TestGetCurrentTime_DateOnly(t *testing.T) {
	ctx := context.Background()
	input := TimeInput{Format: "date"}

	output, err := getCurrentTime(ctx, input)
	if err != nil {
		t.Fatalf("getCurrentTime() error = %v", err)
	}

	// 验证格式 "2006-01-02"（应该只有日期部分，没有空格）
	if strings.Contains(output.Time, " ") {
		t.Errorf("Date format should not contain space, got: %s", output.Time)
	}

	// 验证包含两个横杠
	if strings.Count(output.Time, "-") != 2 {
		t.Errorf("Date format should have 2 dashes, got: %s", output.Time)
	}

	t.Logf("Date: %s", output.Time)
}

// TestGetCurrentTime_TimeOnly 测试仅时间格式
func TestGetCurrentTime_TimeOnly(t *testing.T) {
	ctx := context.Background()
	input := TimeInput{Format: "time"}

	output, err := getCurrentTime(ctx, input)
	if err != nil {
		t.Fatalf("getCurrentTime() error = %v", err)
	}

	// 验证格式 "15:04:05"（应该有两个冒号）
	if strings.Count(output.Time, ":") != 2 {
		t.Errorf("Time format should have 2 colons, got: %s", output.Time)
	}

	t.Logf("Time: %s", output.Time)
}

// TestGetTimezone 测试获取时区信息
func TestGetTimezone(t *testing.T) {
	ctx := context.Background()
	input := map[string]interface{}{}

	output, err := getTimezone(ctx, input)
	if err != nil {
		t.Fatalf("getTimezone() error = %v", err)
	}

	// 验证时区名称不为空
	if output.Timezone == "" {
		t.Error("Expected non-empty timezone name")
	}

	// 验证本地时间和 UTC 时间不为空
	if output.LocalTime == "" {
		t.Error("Expected non-empty local time")
	}
	if output.UTCTime == "" {
		t.Error("Expected non-empty UTC time")
	}

	// 验证时区列表不为空
	if len(output.Timezones) == 0 {
		t.Error("Expected non-empty timezone list")
	}

	// 验证常见时区在列表中
	hasUTC := false
	for _, tz := range output.Timezones {
		if tz == "UTC" {
			hasUTC = true
			break
		}
	}
	if !hasUTC {
		t.Error("Expected UTC in timezone list")
	}

	t.Logf("Timezone: %s, Offset: %d, Local: %s, UTC: %s",
		output.Timezone, output.Offset, output.LocalTime, output.UTCTime)
}

// TestCreateTimeSkill 测试 Skill 创建
func TestCreateTimeSkill(t *testing.T) {
	skill := createTimeSkill()

	if skill == nil {
		t.Fatal("Expected non-nil skill")
	}

	// 验证元数据
	metadata := skill.Glance()
	if metadata == "" {
		t.Error("Expected non-empty metadata")
	}

	// 验证包含时间 skill 的名称
	if !strings.Contains(metadata, "time_skill") {
		t.Errorf("Metadata should contain 'time_skill', got: %s", metadata)
	}

	// 验证 Body 不为空
	body := skill.Inspect()
	if body == "" {
		t.Error("Expected non-empty body")
	}

	// 验证包含 XML 标记
	if !skill.HasXMLTags() {
		t.Error("Expected skill to have XML tags")
	}

	// 验证脚本列表
	scriptNames := skill.GetScriptNames()
	if len(scriptNames) != 2 {
		t.Errorf("Expected 2 scripts, got %d: %v", len(scriptNames), scriptNames)
	}

	// 验证参考文献
	refNames := skill.GetReferenceNames()
	if len(refNames) != 1 {
		t.Errorf("Expected 1 reference, got %d", len(refNames))
	}

	t.Logf("Skill created successfully with scripts: %v", scriptNames)
}

// TestTimeSkill_AutoExecute 测试 Skill 自动执行
func TestTimeSkill_AutoExecute(t *testing.T) {
	skill := createTimeSkill()
	ctx := context.Background()

	input := TimeInput{Format: "local"}
	inputJSON, _ := json.Marshal(input)

	results, err := skill.AutoExecute(ctx, string(inputJSON))
	if err != nil {
		t.Fatalf("AutoExecute() error = %v", err)
	}

	// 验证执行了 2 个脚本
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// 验证每个脚本的执行结果
	for _, result := range results {
		if result.Error != nil {
			t.Errorf("Script %s failed: %v", result.ScriptName, result.Error)
		}
		if result.Result == "" {
			t.Errorf("Script %s returned empty result", result.ScriptName)
		}
		t.Logf("Script %s: %s", result.ScriptName, result.Result)
	}
}

// TestTimeSkill_Execute 测试 Skill 完整执行
func TestTimeSkill_Execute(t *testing.T) {
	skill := createTimeSkill()
	ctx := context.Background()

	input := TimeInput{Format: "iso"}
	inputJSON, _ := json.Marshal(input)

	result, err := skill.Execute(ctx, string(inputJSON))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	t.Logf("Execute result: %s", result)
}

// TestTimeSkill_ReadReference 测试读取参考文献
func TestTimeSkill_ReadReference(t *testing.T) {
	skill := createTimeSkill()

	content, err := skill.ReadReference("time_format_guide")
	if err != nil {
		t.Fatalf("ReadReference() error = %v", err)
	}

	if content == "" {
		t.Error("Expected non-empty reference content")
	}

	// 验证包含时间格式相关内容
	if !strings.Contains(content, "时间") && !strings.Contains(content, "Time") {
		t.Error("Reference content should contain time-related information")
	}

	t.Logf("Reference content length: %d", len(content))
}

// TestGetCurrentTime_DefaultFormat 测试默认格式
func TestGetCurrentTime_DefaultFormat(t *testing.T) {
	ctx := context.Background()
	input := TimeInput{} // 空输入，应该使用默认格式

	output, err := getCurrentTime(ctx, input)
	if err != nil {
		t.Fatalf("getCurrentTime() error = %v", err)
	}

	// 默认应该是 ISO 格式
	if !strings.Contains(output.Time, "T") {
		t.Errorf("Default format should be ISO (contain 'T'), got: %s", output.Time)
	}

	t.Logf("Default format time: %s", output.Time)
}

// BenchmarkGetCurrentTime 性能测试
func BenchmarkGetCurrentTime(b *testing.B) {
	ctx := context.Background()
	input := TimeInput{Format: "iso"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := getCurrentTime(ctx, input)
		if err != nil {
			b.Fatalf("getCurrentTime() error = %v", err)
		}
	}
}
