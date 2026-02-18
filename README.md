# skill

## 项目介绍

`github.com/alois132/skill` 是一个基于 Golang 的 skill 库，专为 AI Agent 应用开发设计。它提供了一套模块化的 skill 系统，让开发者能够轻松构建、管理和部署可在生产环境中使用的 agent skill。

### 核心目标

- **可重用性**：构建可在多个 Agent 应用间共享的 skill 模块
- **类型安全**：利用 Golang 泛型实现编译时类型检查
- **生产就绪**：提供可直接应用于生产项目的完整解决方案
- **框架集成**：深度集成主流 AI Agent 框架（如 CloudWeGo Eino）

### 设计理念

本项目的设计哲学深受 Claude Code 启发，致力于创建一套标准化的 AI Agent 技能开发体系，让开发者能够像构建软件组件一样构建和组合 AI 技能。

## 核心架构

### Skill 结构

Skill 是系统的核心抽象，代表一个 AI Agent 可使用的技能单元。

**属性**：
- `Name`：技能唯一标识
- `Desc`：技能描述
- `Body`：技能核心逻辑描述
- `Scripts`：可执行的脚本列表
- `References`：参考文献/知识库
- `Assets`：多媒体资源（图片、PPT、字体等）

### Resources 体系

Skill 由三类资源构成，形成完整的知识表示体系：

#### 1. Reference（参考文献）
- 文本形式的背景知识
- 提供技能执行的上下文和约束
- 支持 markdown 格式，便于大模型理解

#### 2. Script（可执行脚本）
- 类型安全的可执行函数
- 支持泛型参数，自动序列化/反序列化
- 内置 CLI 工具辅助生成脚本
- **特点**：
  - 运行时类型检查
  - JSON 参数的自动转换
  - 错误处理和日志记录

#### 3. Asset（媒体资产）
- 二进制资源管理
- 支持图片、PPT、字体等多种格式
- 为技能提供视觉辅助材料

### 核心技术

#### 泛型反射系统
借鉴 CloudWeGo 的优秀实践，提供安全的泛型实例创建能力：
- `MakeMap[K, V]()` - 创建映射
- `MakeSlice[T]()` - 创建切片
- `MakePointer[T]()` - 创建指针
- `MakeAny(typeString)` - 基于类型字符串创建实例

#### Builder 模式
流畅的 API 设计：
```go
skill := core.CreateSkill(
    "data_analysis",
    "Perform data analysis and generate insights",
    core.WithScript(analyzeScript),
    core.WithReference(usageRef),
)
```

## 开发计划

### Phase 1：核心功能 ✅
- [x] Skill 核心数据结构定义
- [x] Resources 抽象接口设计（Script、Reference、Asset）
- [x] Builder 模式实现
- [x] 泛型反射工具库

### Phase 2：框架集成（进行中）
- [ ] CloudWeGo Eino 适配器开发
- [ ] Tool 转换层：`{skill_name}()`、`read_reference()`、`use_script()`
- [ ] 中间件和拦截器支持
- [ ] 示例项目：Eino + Skill 集成

### Phase 3：增强功能
- [ ] Asset 功能完善（媒体资源管理）
- [ ] CLI 工具：skill init、skill validate、skill publish
- [ ] Skill 市场/注册中心原型
- [ ] 可视化 skill 编辑器（Web UI）

### Phase 4：生产就绪
- [ ] 完整的测试覆盖率（单元测试 + 集成测试）
- [ ] Benchmark 性能测试
- [ ] 文档和示例完善
- [ ] CI/CD 流水线
- [ ] Docker 部署支持

## XML 标记语法

Skill 支持在 Body 中使用 XML 语法，让大模型能自动识别和调用相关资源。

### 支持的 XML 标记

#### 1. 脚本标记 `<script>`
在 Body 中嵌入脚本引用，自动执行相关脚本。

**语法：**
```go
body := `
第一步：使用<script>init_skill</script>初始化skill
`
```

**完整示例：**
```go
// 定义脚本
func initSkill(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    return map[string]interface{}{
        "message": "Skill initialized",
        "files": []string{"main.go", "README.md"},
    }, nil
}

// 创建 skill
skill := core.CreateSkill(
    "initializer",
    "Initialize a new skill",
    core.WithScript(core.CreateScript("init_skill", initSkill)),
    core.WithAutoParsedBody(`
第一步：使用<script>init_skill</script>初始化skill

此脚本将：
- 创建基础目录结构
- 生成配置文件
- 准备开发环境
`),
)

// 自动执行 Body 中所有脚本
results, err := core.AutoExecute(ctx, skill, `{"skill_name":"demo"}`)
// 或执行完整逻辑
output, err := core.Execute(ctx, skill, `{"skill_name":"demo"}`)
```

#### 2. 参考文献标记 `<reference>`
在 Body 中引用参考文献文档：

```go
body := `
参考文档：<reference>usage_guide</reference>
配置文件：<reference>config_schema</reference>
`

skill := core.CreateSkill(
    "demo",
    "Demo skill",
    core.WithReference("usage_guide", "使用说明内容"),
    core.WithReference("config_schema", "配置说明内容"),
    core.WithBody(body),
)

// 获取参考文献
content, err := core.ReadReference(skill, "usage_guide")
```

#### 3. 资产标记 `<asset>`
在 Body 中引用媒体资产：

```go
body := `
使用模板：<asset>project_template.png</asset>
`

asset := core.CreateAsset("project_template", imageData, resources.PNG)
skill := core.CreateSkill(
    "demo",
    "Demo skill",
    core.WithAsset(asset),
    core.WithBody(body),
)

// 获取资产信息
names := core.GetAssetNames(skill)
```

### Builder 辅助函数

使用辅助函数生成 XML 标记字符串：

```go
body := fmt.Sprintf(`
使用 %s 初始化
参考 %s
模板 %s
`,
    core.EmbedScript("init"),
    core.EmbedReference("guide"),
    core.EmbedAsset("template.png"),
)
```

### 自动执行方法

#### AutoExecute
按顺序执行 Body 中所有 `<script>` 标记对应的脚本：

```go
results, err := skill.AutoExecute(ctx, args)
// results: []ScriptResult{
//     {ScriptName: "init", Result: "...", Error: nil},
//     {ScriptName: "config", Result: "...", Error: nil},
// }
```

#### Execute
执行完整的 skill 逻辑，返回格式化结果：

```go
output, err := skill.Execute(ctx, args)
// 返回格式：
// Skill: skill_name
//
// [1] Script: init_skill
// Result: {...}
//
// [2] Script: config_skill
// Result: {...}
```

### 快速开始

### 安装

```bash
go get github.com/alois132/skill
```

### 基础示例

```go
package main

import (
    "context"
    "github.com/alois132/skill/core"
)

// 定义一个分析脚本
func analyzeData(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    // 分析逻辑
    return map[string]interface{}{
        "result": "analysis completed",
    }, nil
}

func main() {
    // 创建 skill
    analysisSkill := core.CreateSkill(
        "data_analyzer",
        "Analyze data and provide insights",
        core.WithScript(core.CreateScript("analyze", analyzeData)),
        core.WithReference("usage_guide", "Detailed usage instructions..."),
    )

    // 查看 skill 信息
    analysisSkill.Glance()
    analysisSkill.Inspect()
}
```

### 进阶示例：使用 XML 标记

```go
package main

import (
    "context"
    "fmt"
    "github.com/alois132/skill/core"
)

func initProject(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    return map[string]interface{}{
        "message": "Project initialized",
    }, nil
}

func setupConfig(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    return map[string]interface{}{
        "config": "Configuration applied",
    }, nil
}

func main() {
    skill := core.CreateSkill(
        "project_setup",
        "Set up a new project",
        core.WithScript(core.CreateScript("init", initProject)),
        core.WithScript(core.CreateScript("config", setupConfig)),
        core.WithAutoParsedBody(`
第一步：<script>init</script> - 初始化项目
第二步：<script>config</script> - 配置文件
参考：<reference>setup_guide</reference>
`),
    )

    // 自动执行所有脚本
    results, err := core.AutoExecute(context.Background(), skill, `{}`)
    if err != nil {
        panic(err)
    }

    for _, r := range results {
        fmt.Printf("Script: %s\nResult: %s\n\n", r.ScriptName, r.Result)
    }
}
```

## 扩展：与 CloudWeGo Eino 集成

skill 项目专门设计与 CloudWeGo Eino 框架无缝集成，将 Skill 封装为 Eino Tool。

### 三类 Tool 映射

1. **`{skill_name}()`**
   - **用途**：执行完整的 skill 逻辑
   - **示例**：创建 `skill_create` skill 后，自动生成 `skill_create()` tool
   - **输入**：skill 的 body 作为描述
   - **输出**：skill 执行结果

2. **`read_reference(reference_name string)`**
   - **用途**：访问 skill 的参考文献
   - **示例**：`read_reference("skill_create/workflows")`
   - **输入**：reference 名称（格式：`skill_name/reference_name`）
   - **输出**：reference 的 body 内容

3. **`use_script(script_name string, args string)`**
   - **用途**：执行 skill 中的特定脚本
   - **示例**：`use_script("skill_create/init_skill", `{"skill_name":"skill_create"}`)
   - **输入**：脚本名称和 JSON 格式的参数
   - **输出**：脚本执行结果的 JSON 字符串

### 集成优势

- ✅ 统一的 Tool 接口标准
- ✅ 类型安全的参数处理
- ✅ 自动化的 skill 发现机制
- ✅ 与 Eino 生态无缝协作

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

### 开发环境准备

```bash
# 克隆项目
git clone https://github.com/alois132/skill.git
cd skill

# 安装依赖
go mod download

# 运行测试
go test ./...
```

### 代码规范

- 遵循 Go 官方编码规范
- 中英双语注释（代码用英文，说明用中文）
- 提交前确保 `go vet` 和 `go fmt` 通过