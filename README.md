# Skill

<p align="center">
  <strong>A modular skill framework for AI Agent applications in Go</strong><br>
  <strong>基于 Go 的 AI Agent 模块化技能框架</strong>
</p>

<p align="center">
  <a href="https://github.com/alois132/skill/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  </a>
  <a href="https://pkg.go.dev/github.com/alois132/skill">
    <img src="https://pkg.go.dev/badge/github.com/alois132/skill.svg" alt="Go Reference">
  </a>
  <a href="https://goreportcard.com/report/github.com/alois132/skill">
    <img src="https://goreportcard.com/badge/github.com/alois132/skill" alt="Go Report Card">
  </a>
  <a href="https://github.com/alois132/skill/actions">
    <img src="https://img.shields.io/badge/tests-passing-brightgreen.svg" alt="Tests">
  </a>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#documentation">Documentation</a> •
  <a href="#eino-integration">Eino Integration</a> •
  <a href="#examples">Examples</a>
</p>

---

## Overview

**Skill** is a production-ready Go library for building modular, type-safe, and reusable AI Agent skills. It provides a standardized way to define agent capabilities with scripts, references, and assets, while offering seamless integration with popular AI frameworks like [CloudWeGo Eino](https://github.com/cloudwego/eino).

### Key Features

- **Modular Design** - Build reusable skill modules that can be shared across multiple Agent applications
- **Type Safety** - Leverage Go generics for compile-time type checking and automatic JSON serialization
- **XML Tag Syntax** - Use intuitive XML tags (`<script>`, `<reference>`, `<asset>`) in skill body for AI to identify resources
- **Framework Integration** - Deep integration with CloudWeGo Eino, converting skills to tools automatically
- **Production Ready** - Complete with error handling, logging, and comprehensive testing

### Use Cases

- 🤖 **AI Agent Development** - Build intelligent agents with structured capabilities
- 🔧 **Tool Orchestration** - Organize and execute multiple tools in a coordinated manner
- 📚 **Knowledge Management** - Manage references and documentation for AI context
- 🎨 **Multimedia Skills** - Handle assets like images, templates, and documents

---

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Core Concepts](#core-concepts)
  - [Skill Structure](#skill-structure)
  - [Resources](#resources)
  - [XML Tag Syntax](#xml-tag-syntax)
- [API Reference](#api-reference)
  - [Creating Skills](#creating-skills)
  - [Executing Scripts](#executing-scripts)
  - [Reading References](#reading-references)
- [Eino Integration](#eino-integration)
- [Examples](#examples)
- [Development Roadmap](#development-roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## Installation

```bash
go get github.com/alois132/skill
```

### Requirements

- Go 1.18 or higher (for generics support)

---

## Quick Start

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/alois132/skill/core"
)

// Define a script function
func analyzeData(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    return map[string]interface{}{
        "result": "Analysis completed successfully",
        "data":   input,
    }, nil
}

func main() {
    // Create a skill
    skill := core.CreateSkill(
        "data_analyzer",
        "Analyze data and provide insights",
        core.WithScript(core.CreateScript("analyze", analyzeData)),
        core.WithReference("usage_guide", "Detailed usage instructions..."),
    )

    // Use a specific script from the skill
    ctx := context.Background()
    result, err := skill.UseScript(ctx, "analyze", `{"query": "sample data"}`)
    if err != nil {
        panic(err)
    }
    fmt.Println(result)
}
```

### XML Tag Example

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
        "files":   []string{"main.go", "README.md"},
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
        "Set up a new Go project",
        core.WithScript(core.CreateScript("init", initProject)),
        core.WithScript(core.CreateScript("config", setupConfig)),
        core.WithAutoParsedBody(`
第一步：<script>init</script> - 初始化项目结构
第二步：<script>config</script> - 配置项目参数

参考文档：<reference>setup_guide</reference>
        `),
    )

    // Use a specific script from the skill
    ctx := context.Background()
    result, err := skill.UseScript(ctx, "init", `{"project_name":"myapp"}`)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Result: %v\n", result)
}
```

---

## Core Concepts

### Skill Structure

A `Skill` is the core abstraction representing a capability that an AI Agent can use:

| Field | Type | Description |
|-------|------|-------------|
| `Name` | `string` | Unique identifier for the skill |
| `Desc` | `string` | Human-readable description |
| `Body` | `string` | Core logic description with XML tags |
| `Scripts` | `[]Script` | Executable functions |
| `References` | `[]Reference` | Documentation and knowledge base |
| `Assets` | `[]Asset` | Binary resources (images, templates, etc.) |

### Resources

Skills are composed of three types of resources:

#### 1. Script - Executable Functions

Type-safe executable functions with automatic JSON serialization:

```go
// Define with generics for type safety
script := resources.NewEasyScript(
    "calculate",
    func(ctx context.Context, req CalculateRequest) (CalculateResponse, error) {
        return CalculateResponse{Result: req.A + req.B}, nil
    },
)
```

Features:
- Runtime type checking
- Automatic JSON parameter conversion
- Generic type safety

#### 2. Reference - Documentation

Text-based knowledge resources:

```go
ref := resources.Reference{
    Name: "api_guide",
    Body: "# API Documentation\n...",
}
```

#### 3. Asset - Binary Resources

Binary media resources:

```go
asset := resources.Asset{
    Name: "template",
    Type: resources.PNG,
    Body: imageData,
}
```

### XML Tag Syntax

Use XML tags in the skill body to help AI identify and use resources:

| Tag | Purpose | Example |
|-----|---------|---------|
| `<script>name</script>` | Reference a script | `<script>init_project</script>` |
| `<reference>name</reference>` | Reference documentation | `<reference>api_guide</reference>` |
| `<asset>name</asset>` | Reference a binary asset | `<asset>template.png</asset>` |

---

## API Reference

### Creating Skills

```go
// Basic skill creation
skill := core.CreateSkill(
    "skill_name",
    "Skill description",
    core.WithScript(script),
    core.WithReference("guide", "Documentation content"),
    core.WithBody("Skill body with <script>example</script>"),
)

// With auto-parsed body
skill := core.CreateSkill(
    "skill_name",
    "Skill description",
    core.WithAutoParsedBody(`
步骤：<script>step1</script>
参考：<reference>guide</reference>
    `),
)
```

### Executing Scripts

```go
// Execute a specific script
result, err := skill.UseScript(ctx, "script_name", `{"key":"value"}`)
```

### Reading References

```go
// Read a reference
content, err := skill.ReadReference("reference_name")

// Get all reference names
names := skill.GetReferenceNames()
```

---

## Eino Integration

Skill provides seamless integration with [CloudWeGo Eino](https://github.com/cloudwego/eino), automatically converting skills to Eino Tools.

### Converting Skills to Tools

```go
import "github.com/alois132/skill/adapter/eino"

// Convert a skill to Eino tools
tools, err := eino.ToTools(skill)

// Convert multiple skills
allTools, err := eino.ToTools(skill1, skill2, skill3...)
```

### Generated Tools

Each skill generates tools for accessing its resources:

| Tool | Purpose | Input | Output |
|------|---------|-------|--------|
| `read_reference(name)` | Access documentation | Reference name | Reference content |
| `use_script(name, args)` | Execute specific script | Script name + JSON args | Script result |

### Example: Eino Agent with Skills

```go
package main

import (
    "context"
    "github.com/alois132/skill/adapter/eino"
    "github.com/alois132/skill/core"
    "github.com/cloudwego/eino/components/tool"
)

func main() {
    // Create skills
    timeSkill := createTimeSkill()
    calcSkill := createCalcSkill()

    // Convert to Eino tools
    tools, err := eino.ToTools(timeSkill, calcSkill)
    if err != nil {
        panic(err)
    }

    // Use with Eino agent
    agent := createAgent(tools)
    response, err := agent.Generate(ctx, "What time is it?")
}
```

---

## Examples

### Time Skill Example

A complete example demonstrating time formatting capabilities:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/alois132/skill/core"
    "github.com/alois132/skill/schema/resources"
)

func main() {
    // Create time formatting skill
    timeSkill := core.CreateSkill(
        "time_formatter",
        "Format and manipulate time",
        core.WithScript(resources.NewEasyScript(
            "format_time",
            func(ctx context.Context, req struct {
                Format string `json:"format"`
            }) (struct {
                Result string `json:"result"`
            }, error) {
                return struct {
                    Result string `json:"result"`
                }{
                    Result: time.Now().Format(req.Format),
                }, nil
            },
        )),
        core.WithReference("time_formats", `
常用时间格式：
- "2006-01-02" - 日期
- "15:04:05" - 时间
- "2006-01-02 15:04:05" - 完整时间
        `),
        core.WithAutoParsedBody(`
格式化当前时间：<script>format_time</script>
参考时间格式文档：<reference>time_formats</reference>
        `),
    )

    // Execute
    result, err := timeSkill.UseScript(context.Background(), "format_time", `{"format":"2006-01-02 15:04:05"}`)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Result: %v\n", result)
}
```

### More Examples

See the [`example/`](./example) directory for more complete examples:

- [`main.go`](./example/main.go) - Basic usage demonstration
- [`time_skill.go`](./example/time_skill.go) - Time formatting skill with tests
- [`init_skill.go`](./example/init_skill.go) - Project initialization with XML tags

---

## Development Roadmap

### Phase 1: Core Foundation ✅
- [x] Skill core data structures
- [x] Resource abstractions (Script, Reference, Asset)
- [x] Builder pattern API
- [x] Generic reflection utilities

### Phase 2: Framework Integration ✅
- [x] CloudWeGo Eino adapter
- [x] Tool conversion layer
- [x] Automatic skill discovery

### Phase 3: Enhanced Features 🚧
- [ ] Asset management system
- [ ] CLI tools (skill init, validate, publish)
- [ ] Skill registry/marketplace
- [ ] Web-based skill editor

### Phase 4: Production Readiness
- [ ] Comprehensive test coverage
- [ ] Performance benchmarks
- [ ] CI/CD pipeline
- [ ] Docker deployment support

---

## Contributing

We welcome contributions! Please see our [contributing guidelines](./CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/alois132/skill.git
cd skill

# Install dependencies
go mod download

# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

### Code Standards

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use bilingual comments (English for code, Chinese for explanations)
- Ensure `go vet` and `go fmt` pass before committing

---

## Related Projects

- [CloudWeGo Eino](https://github.com/cloudwego/eino) - The AI Agent framework that Skill integrates with
- [LangChain Go](https://github.com/tmc/langchaingo) - Another popular Go AI framework
- [Claude Code](https://github.com/anthropics/claude-code) - Inspiration for this project's design philosophy

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  <strong>Built with ❤️ for the AI Agent community</strong><br>
  <sub>Designed for production, inspired by Claude Code</sub>
</p>
