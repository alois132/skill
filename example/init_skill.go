package main

import (
	"context"
	"fmt"

	"github.com/alois132/skill/core"
)

// 初始化脚本：创建基础 skill 结构
func initSkill(_ context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	skillName, ok := input["skill_name"].(string)
	if !ok || skillName == "" {
		return nil, fmt.Errorf("无效的 skill_name: %v", input["skill_name"])
	}

	// 模拟初始化操作
	result := map[string]interface{}{
		"message": fmt.Sprintf("Skill '%s' initialized successfully", skillName),
		"files_created": []string{
			fmt.Sprintf("%s/main.go", skillName),
			fmt.Sprintf("%s/README.md", skillName),
			fmt.Sprintf("%s/scripts/init.py", skillName),
			fmt.Sprintf("%s/docs/usage.md", skillName),
		},
		"directories_created": []string{
			fmt.Sprintf("%s/", skillName),
			fmt.Sprintf("%s/scripts/", skillName),
			fmt.Sprintf("%s/docs/", skillName),
		},
	}

	return result, nil
}

// 配置脚本：生成配置文件
func configSkill(_ context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	skillName := input["skill_name"].(string)
	configType := "default"
	if ct, ok := input["config_type"].(string); ok {
		configType = ct
	}

	result := map[string]interface{}{
		"message": fmt.Sprintf("Config generated for %s with type: %s", skillName, configType),
		"config_file": fmt.Sprintf("%s/config/%s.json", skillName, configType),
	}

	return result, nil
}

func main() {
	fmt.Println("=== Skill Initialization Demo with XML Syntax ===\n")

	// 创建初始化 skill，使用 XML 标记语法
	initSkill := core.CreateSkill(
		"skill_init",
		"Initialize a new skill with directory structure and configuration",
		core.WithScript(core.CreateScript("init_skill", initSkill)),
		core.WithScript(core.CreateScript("config_skill", configSkill)),
		core.WithAutoParsedBody(`
第一步：使用<script>init_skill</script>初始化skill

此脚本将：
- 创建基础目录结构
- 生成示例文件
- 准备开发环境

第二步：使用<script>config_skill</script>生成配置文件

可以指定配置类型：
- default: 默认配置
- dev: 开发环境配置
- prod: 生产环境配置

参考文档：<reference>usage_guide</reference>
使用模板：<asset>project_template.png</asset>
`),
		core.WithReference("usage_guide", `# 使用指南

## 初始化步骤
1. 运行 init_skill 脚本创建项目结构
2. 运行 config_skill 生成配置文件
3. 开始开发

## 目录结构
- main.go: 主入口
- scripts/: 脚本目录
- docs/: 文档目录
`),
	)

	// 1. 查看 Skill 基本信息
	fmt.Println("1. Skill Metadata:")
	metadata := initSkill.Glance()
	fmt.Println(metadata)

	// 2. 检查 XML 标记解析
	fmt.Println("\n2. XML Tags in Body:")
	if initSkill.HasXMLTags() {
		fmt.Println("✓ Body contains XML tags")
		scriptNames := initSkill.GetScriptNames()
		fmt.Printf("  - Scripts found: %v\n", scriptNames)

		refNames := initSkill.GetReferenceNames()
		fmt.Printf("  - References found: %v\n", refNames)

		assetNames := initSkill.GetAssetNames()
		fmt.Printf("  - Assets found: %v\n", assetNames)
	} else {
		fmt.Println("✗ No XML tags found")
	}

	// 3. 演示嵌入函数
	fmt.Println("\n5. Builder Helper Functions:")
	body := fmt.Sprintf(`
使用 %s 初始化
参考 %s
模板 %s
`,
		core.EmbedScript("init"),
		core.EmbedReference("guide"),
		core.EmbedAsset("template.png"),
	)
	fmt.Println(body)

	fmt.Println("\n=== Demo completed successfully! ===")
}
