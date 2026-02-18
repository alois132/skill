package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alois132/skill/core"
	"github.com/alois132/skill/schema/resources"
)

// 示例1：基础数据分析脚本
func analyzeData(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	// 模拟数据分析
	result := map[string]interface{}{
		"message": "analysis completed",
		"data":    input["data"],
	}
	return result, nil
}

// 示例2：字符串处理脚本
func processString(ctx context.Context, input string) (string, error) {
	return "processed: " + input, nil
}

func main() {
	fmt.Println("=== Skill Library Demo ===\n")

	// 创建data_analysis skill
	dataAnalysisSkill := core.CreateSkill(
		"data_analyzer",
		"Analyze data and provide insights",
		core.WithScript(core.CreateScript("analyze", analyzeData)),
		core.WithReference("usage_guide", "Detailed usage instructions for data analysis skill..."),
	)

	// Glance查看skill的元数据
	fmt.Println("1. Glance skill metadata:")
	metadata := dataAnalysisSkill.Glance()
	fmt.Println(metadata)

	// Inspect查看skill的详细信息
	fmt.Println("\n2. Inspect skill body:")
	// 当前body是空的，因为我们没有设置
	body := dataAnalysisSkill.Inspect()
	fmt.Println(body)

	// 使用UseScript执行脚本
	fmt.Println("\n3. Use skill script:")
	result, err := core.UseScript(
		context.Background(),
		dataAnalysisSkill,
		"analyze",
		`{"data": "sample data"}`,
	)
	if err != nil {
		log.Printf("Error using script: %v", err)
	} else {
		fmt.Printf("Script result: %s\n", result)
	}

	// 使用ReadReference读取参考文献
	fmt.Println("\n4. Read skill reference:")
	ref, err := core.ReadReference(dataAnalysisSkill, "usage_guide")
	if err != nil {
		log.Printf("Error reading reference: %v", err)
	} else {
		fmt.Printf("Reference content: %s\n", ref)
	}

	// 创建一个更复杂的skill示例
	fmt.Println("\n=== Advanced Example ===")

	// 创建多个脚本
	stringScript := core.CreateScript("process", processString)

	// 创建多个参考文献
	ref1 := core.CreateReference("best_practices", "# Best Practices\n1. Validate input\n2. Handle errors\n3. Log results")
	ref2 := core.CreateReference("api_docs", "# API Documentation\nThis API provides...")

	// 创建资产
	sampleData := []byte("sample binary data")
	asset := core.CreateAsset("sample_template", sampleData, resources.PNG)

	// 使用Builder模式创建skill
	advancedSkill := core.CreateSkill(
		"advanced_analyzer",
		"Advanced data analysis with templates and best practices",
		core.WithReferences([]*resources.Reference{ref1, ref2}),
		core.WithScript(stringScript),
		core.WithAsset(asset),
		core.WithBody("This skill provides advanced analysis capabilities..."),
	)

	// 测试新的skill
	metadata2 := advancedSkill.Glance()
	fmt.Println("Advanced skill metadata:")
	fmt.Println(metadata2)

	// 测试读取best_practices参考文献
	refContent, err := core.ReadReference(advancedSkill, "best_practices")
	if err != nil {
		log.Printf("Error reading reference: %v", err)
	} else {
		refObj := resources.Reference{Name: "best_practices", Body: refContent}
		fmt.Printf("\nBest practices (summary): %s\n", refObj.Summary())
	}

	fmt.Println("\n=== Demo completed successfully! ===")
}
