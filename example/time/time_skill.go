package time

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alois132/skill/core"
	"github.com/alois132/skill/schema"
)

// TimeInput 时间脚本输入参数
type TimeInput struct {
	Format   string `json:"format"`   // 时间格式: iso, local, unix, custom
	Timezone string `json:"timezone"` // 时区，如 "UTC", "Asia/Shanghai"
	Layout   string `json:"layout"`   // 自定义格式布局（当 format=custom 时使用）
}

// TimeOutput 时间脚本输出结果
type TimeOutput struct {
	Time     string `json:"time"`     // 格式化后的时间字符串
	Unix     int64  `json:"unix"`     // Unix 时间戳
	Timezone string `json:"timezone"` // 当前时区
}

// TimezoneOutput 时区信息输出
type TimezoneOutput struct {
	Timezone  string   `json:"timezone"`   // 当前时区名称
	Offset    int      `json:"offset"`     // 偏移量（秒）
	LocalTime string   `json:"local_time"` // 本地时间
	UTCTime   string   `json:"utc_time"`   // UTC时间
	Timezones []string `json:"timezones"`  // 常见时区列表
}

// getCurrentTime 获取当前时间
// 支持多种格式：iso, local, unix, custom
func getCurrentTime(ctx context.Context, input TimeInput) (TimeOutput, error) {
	// 确定时区
	loc := time.Local
	if input.Timezone != "" {
		var err error
		loc, err = time.LoadLocation(input.Timezone)
		if err != nil {
			return TimeOutput{}, fmt.Errorf("invalid timezone: %s", input.Timezone)
		}
	}

	now := time.Now().In(loc)
	output := TimeOutput{
		Unix:     now.Unix(),
		Timezone: loc.String(),
	}

	// 根据格式要求格式化时间
	switch input.Format {
	case "unix":
		output.Time = fmt.Sprintf("%d", now.Unix())
	case "iso", "":
		output.Time = now.Format(time.RFC3339)
	case "local":
		output.Time = now.Format("2006-01-02 15:04:05")
	case "date":
		output.Time = now.Format("2006-01-02")
	case "time":
		output.Time = now.Format("15:04:05")
	case "custom":
		if input.Layout != "" {
			output.Time = now.Format(input.Layout)
		} else {
			output.Time = now.Format("2006-01-02 15:04:05")
		}
	default:
		output.Time = now.Format(time.RFC3339)
	}

	return output, nil
}

// getTimezone 获取时区信息
func getTimezone(ctx context.Context, input map[string]interface{}) (TimezoneOutput, error) {
	now := time.Now()
	_, offset := now.Zone()

	output := TimezoneOutput{
		Timezone:  time.Local.String(),
		Offset:    offset,
		LocalTime: now.Format("2006-01-02 15:04:05"),
		UTCTime:   now.UTC().Format("2006-01-02 15:04:05"),
		Timezones: []string{
			"UTC",
			"Asia/Shanghai",
			"Asia/Tokyo",
			"Asia/Seoul",
			"Asia/Singapore",
			"Europe/London",
			"Europe/Paris",
			"Europe/Berlin",
			"America/New_York",
			"America/Los_Angeles",
			"America/Chicago",
			"Australia/Sydney",
		},
	}

	return output, nil
}

// createTimeSkill 创建时间 Skill
func createTimeSkill() *schema.Skill {
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
  - "local" - 本地时间格式 (2006-01-02 15:04:05)
  - "unix" - Unix 时间戳
  - "date" - 仅日期
  - "time" - 仅时间
  - "custom" - 自定义格式（配合 layout 参数）
- timezone: 时区名称，如 "UTC", "Asia/Shanghai"
- layout: 自定义格式布局（当 format="custom" 时使用）

### 2. 获取时区信息
使用 <script>get_timezone</script> 脚本获取当前时区信息和常见时区列表。

## 示例

获取本地时间：
{"format": "local"}

获取 ISO 格式时间：
{"format": "iso"}

获取上海时间：
{"format": "local", "timezone": "Asia/Shanghai"}

获取 Unix 时间戳：
{"format": "unix"}

自定义格式：
{"format": "custom", "layout": "2006年01月02日 15:04"}

## 参考文档
更多时间格式说明请参考：<reference>time_format_guide</reference>
`),
		core.WithReference("time_format_guide", `# 时间格式指南

## Go 时间布局参考
Go 使用以下参考时间进行格式化：
  "Mon Jan 2 15:04:05 MST 2006"

常用格式：
- "2006-01-02"                    - 日期
- "15:04:05"                      - 时间
- "2006-01-02 15:04:05"          - 日期时间
- "2006年01月02日"               - 中文日期
- "Jan 2, 2006"                  - 英文日期
- "Monday, January 2, 2006"      - 完整英文日期

## 时区参考
常用时区：
- UTC              - 协调世界时
- Asia/Shanghai    - 中国标准时间 (CST)
- Asia/Tokyo       - 日本标准时间 (JST)
- Asia/Seoul       - 韩国标准时间 (KST)
- Europe/London    - 格林尼治时间 (GMT)
- Europe/Paris     - 中欧时间 (CET)
- America/New_York - 美国东部时间 (EST)
- America/Los_Angeles - 美国太平洋时间 (PST)
`),
	)
}

// RunTimeSkillDemo 运行时间 Skill 演示
func RunTimeSkillDemo() {
	fmt.Println("=== Time Skill Demo ===")
	fmt.Println()

	// 创建时间 Skill
	timeSkill := createTimeSkill()

	// 1. 查看 Skill 基本信息
	fmt.Println("1. Skill Metadata:")
	metadata := timeSkill.Glance()
	fmt.Println(metadata)

	// 2. 检查 XML 标记解析
	fmt.Println("\n2. XML Tags in Body:")
	if timeSkill.HasXMLTags() {
		fmt.Println("✓ Body contains XML tags")
		scriptNames := timeSkill.GetScriptNames()
		fmt.Printf("  - Scripts found: %v\n", scriptNames)

		refNames := timeSkill.GetReferenceNames()
		fmt.Printf("  - References found: %v\n", refNames)
	}

	// 3. 执行获取时间脚本
	fmt.Println("\n3. Get Current Time (ISO format):")
	ctx := context.Background()

	isoInput := TimeInput{Format: "iso"}
	isoJSON, _ := json.Marshal(isoInput)
	result, err := core.UseScript(ctx, timeSkill, "get_current_time", string(isoJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", result)
	}

	// 4. 获取上海时间
	fmt.Println("\n4. Get Shanghai Time:")
	shInput := TimeInput{Format: "local", Timezone: "Asia/Shanghai"}
	shJSON, _ := json.Marshal(shInput)
	result, err = core.UseScript(ctx, timeSkill, "get_current_time", string(shJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", result)
	}

	// 5. 获取 Unix 时间戳
	fmt.Println("\n5. Get Unix Timestamp:")
	unixInput := TimeInput{Format: "unix"}
	unixJSON, _ := json.Marshal(unixInput)
	result, err = core.UseScript(ctx, timeSkill, "get_current_time", string(unixJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", result)
	}

	// 6. 获取时区信息
	fmt.Println("\n6. Get Timezone Info:")
	tzResult, err := core.UseScript(ctx, timeSkill, "get_timezone", "{}")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", tzResult)
	}

	// 7. 读取参考文献
	fmt.Println("\n7. Read Reference:")
	ref, err := core.ReadReference(timeSkill, "time_format_guide")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Reference content:\n%s\n", ref)
	}

	fmt.Println("\n=== Demo completed successfully! ===")
}
