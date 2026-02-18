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

## 快速开始

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
    "github.com/alois132/skill/schema"
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
        core.WithScript(schema.NewEasyScript("analyze", analyzeData)),
        core.WithReference("usage_guide", "Detailed usage instructions..."),
    )

    // 查看 skill 信息
    analysisSkill.Glance()
    analysisSkill.Inspect()
}
```

### 进阶示例：使用 Reference 和 Asset

```go
// 添加参考文献
ref := schema.Reference{
    Name: "best_practices",
    Body: `# 数据分析最佳实践
1. 数据清洗
2. 异常值处理
3. 可视化展示`,
}

// 添加资产
asset := schema.Asset{
    Name: "chart_template",
    Type: "png",
    Data: chartData, // 图表模板数据
}

skill := core.CreateSkill(
    "advanced_analyzer",
    "Advanced data analysis with templates",
    core.WithReference("best_practices", ref),
    core.WithAsset("chart_template", asset),
)
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