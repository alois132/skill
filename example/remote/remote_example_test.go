package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alois132/skill/core"
	"github.com/alois132/skill/schema"
	"github.com/alois132/skill/schema/resources"
)

// TestRemoteScriptExample 演示如何使用 RemoteScript
func TestRemoteScriptExample(t *testing.T) {
	// 创建模拟的远程服务
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req resources.ScriptCallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch req.ScriptName {
		case "weather":
			var args struct {
				City string `json:"city"`
			}
			json.Unmarshal([]byte(req.Args), &args)
			result := fmt.Sprintf(`{"city":"%s","temperature":25,"condition":"Sunny"}`, args.City)
			json.NewEncoder(w).Encode(resources.ScriptCallResponse{Result: result})

		case "translate":
			var args struct {
				Text   string `json:"text"`
				Target string `json:"target"`
			}
			json.Unmarshal([]byte(req.Args), &args)
			result := fmt.Sprintf(`{"original":"%s","translated":"[Translated to %s] %s"}`, args.Text, args.Target, args.Text)
			json.NewEncoder(w).Encode(resources.ScriptCallResponse{Result: result})

		default:
			json.NewEncoder(w).Encode(resources.ScriptCallResponse{Error: "unknown script"})
		}
	}))
	defer server.Close()

	// 创建 HTTP 客户端
	httpClient := resources.NewHTTPRemoteScriptClient(
		server.URL,
		resources.WithHeader("X-API-Key", "test-key"),
	)

	// 创建远程脚本
	weatherScript := resources.NewRemoteScript("weather", httpClient).
		WithUsage("Get weather information for a city")
	translateScript := resources.NewRemoteScript("translate", httpClient).
		WithUsage("Translate text to target language")

	// 创建使用远程脚本的 Skill
	skill := core.CreateSkill(
		"remote_services",
		"Skill using remote services",
		core.WithBody(`
Remote Services Skill

Get weather: <script>weather</script>
Translate text: <script>translate</script>
`),
		core.WithScript(weatherScript),
		core.WithScript(translateScript),
	)

	ctx := context.Background()

	// 执行远程天气脚本
	result, err := skill.UseScript(ctx, "weather", `{"city":"Beijing"}`)
	if err != nil {
		t.Fatalf("Failed to use weather script: %v", err)
	}
	fmt.Printf("Weather: %s\n", result)

	// 执行远程翻译脚本
	result, err = skill.UseScript(ctx, "translate", `{"text":"Hello","target":"Chinese"}`)
	if err != nil {
		t.Fatalf("Failed to use translate script: %v", err)
	}
	fmt.Printf("Translation: %s\n", result)
}

// TestMixedLocalRemoteScripts 演示混合使用本地和远程脚本
func TestMixedLocalRemoteScripts(t *testing.T) {
	// 创建模拟的远程服务
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req resources.ScriptCallRequest
		json.NewDecoder(r.Body).Decode(&req)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if req.ScriptName == "remote_compute" {
			var args struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			}
			json.Unmarshal([]byte(req.Args), &args)
			result := fmt.Sprintf(`{"result":%f,"source":"remote"}`, args.X*args.Y)
			json.NewEncoder(w).Encode(resources.ScriptCallResponse{Result: result})
		}
	}))
	defer server.Close()

	// 创建 HTTP 客户端
	httpClient := resources.NewHTTPRemoteScriptClient(server.URL)

	// 创建 Skill，混合本地和远程脚本
	skill := core.CreateSkill(
		"calculator",
		"Calculator with local and remote operations",
		core.WithBody(`
Calculator Skill

Local operations:
- <script>add</script>: Add two numbers locally
- <script>subtract</script>: Subtract two numbers locally

Remote operations:
- <script>remote_compute</script>: Complex computation on remote server
`),
		// 本地脚本
		core.WithScript(core.CreateScript("add", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
			a, _ := input["a"].(float64)
			b, _ := input["b"].(float64)
			return map[string]interface{}{"result": a + b, "source": "local"}, nil
		})),
		core.WithScript(core.CreateScript("subtract", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
			a, _ := input["a"].(float64)
			b, _ := input["b"].(float64)
			return map[string]interface{}{"result": a - b, "source": "local"}, nil
		})),
		// 远程脚本
		core.WithScript(resources.NewRemoteScript("remote_compute", httpClient)),
	)

	ctx := context.Background()

	// 执行本地脚本
	result, _ := skill.UseScript(ctx, "add", `{"a":10,"b":5}`)
	fmt.Printf("Local add: %s\n", result)

	result, _ = skill.UseScript(ctx, "subtract", `{"a":10,"b":5}`)
	fmt.Printf("Local subtract: %s\n", result)

	// 执行远程脚本
	result, _ = skill.UseScript(ctx, "remote_compute", `{"x":10,"y":5}`)
	fmt.Printf("Remote compute: %s\n", result)
}

// TestResourceProviderWithRemoteScripts 演示使用 ResourceProvider 管理远程脚本
func TestResourceProviderWithRemoteScripts(t *testing.T) {
	// 创建模拟的远程服务
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req resources.ScriptCallRequest
		json.NewDecoder(r.Body).Decode(&req)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch req.ScriptName {
		case "service_a":
			json.NewEncoder(w).Encode(resources.ScriptCallResponse{Result: `{"service":"A","status":"ok"}`})
		case "service_b":
			json.NewEncoder(w).Encode(resources.ScriptCallResponse{Result: `{"service":"B","status":"ok"}`})
		}
	}))
	defer server.Close()

	// 创建 HTTP 客户端
	httpClient := resources.NewHTTPRemoteScriptClient(server.URL)

	// 创建资源提供者，包含远程脚本
	provider := core.CreateInlineProvider()
	provider.AddScript(resources.NewRemoteScript("service_a", httpClient))
	provider.AddScript(resources.NewRemoteScript("service_b", httpClient))
	provider.AddReference(&resources.Reference{
		Name: "service_docs",
		Body: "Documentation for remote services",
	})

	// 创建使用 Provider 的 Skill
	skill := core.CreateSkill(
		"multi_service",
		"Skill accessing multiple remote services",
		core.WithBody(`
Multi-Service Skill

Services:
- <script>service_a</script>
- <script>service_b</script>

Docs: <reference>service_docs</reference>
`),
		core.WithResourceProvider(provider),
	)

	ctx := context.Background()

	// 通过 Provider 执行远程脚本
	result, _ := skill.UseScript(ctx, "service_a", `{}`)
	fmt.Printf("Service A: %s\n", result)

	result, _ = skill.UseScript(ctx, "service_b", `{}`)
	fmt.Printf("Service B: %s\n", result)

	// 读取 Provider 中的参考文档
	ref, _ := skill.ReadReference("service_docs")
	fmt.Printf("Docs: %s\n", ref)
}

// TestMockRemoteClient 演示使用 Mock 客户端进行测试
func TestMockRemoteClient(t *testing.T) {
	// 创建 Mock 客户端
	mockClient := resources.NewMockRemoteScriptClient()

	// 注册模拟的远程脚本
	mockClient.Register("mock_weather", func(ctx context.Context, args string) (string, error) {
		var input struct {
			City string `json:"city"`
		}
		json.Unmarshal([]byte(args), &input)
		return fmt.Sprintf(`{"city":"%s","temp":20,"mock":true}`, input.City), nil
	})

	mockClient.Register("mock_time", func(ctx context.Context, args string) (string, error) {
		return `{"time":"2024-01-01 12:00:00","mock":true}`, nil
	})

	// 创建使用 Mock 客户端的 Skill
	skill := core.CreateSkill(
		"mock_services",
		"Skill with mock remote services for testing",
		core.WithBody(`
Mock Services

- <script>mock_weather</script>
- <script>mock_time</script>
`),
		core.WithScript(resources.NewRemoteScript("mock_weather", mockClient)),
		core.WithScript(resources.NewRemoteScript("mock_time", mockClient)),
	)

	ctx := context.Background()

	// 执行 Mock 脚本
	result, _ := skill.UseScript(ctx, "mock_weather", `{"city":"Shanghai"}`)
	fmt.Printf("Mock weather: %s\n", result)

	result, _ = skill.UseScript(ctx, "mock_time", `{}`)
	fmt.Printf("Mock time: %s\n", result)
}

// TestCompositeProviderWithRemote 演示使用 CompositeProvider 组合本地和远程资源
func TestCompositeProviderWithRemote(t *testing.T) {
	// 创建模拟的远程服务
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req resources.ScriptCallRequest
		json.NewDecoder(r.Body).Decode(&req)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resources.ScriptCallResponse{Result: `{"source":"remote"}`})
	}))
	defer server.Close()

	// 创建本地 Provider
	localProvider := core.CreateInlineProvider()
	localProvider.AddScript(core.CreateScript("local_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"source": "local"}, nil
	}))
	localProvider.AddReference(&resources.Reference{Name: "local_ref", Body: "Local reference"})

	// 创建远程 Provider
	remoteProvider := core.CreateInlineProvider()
	remoteProvider.AddScript(resources.NewRemoteScript("remote_script", resources.NewHTTPRemoteScriptClient(server.URL)))
	remoteProvider.AddReference(&resources.Reference{Name: "remote_ref", Body: "Remote reference"})

	// 创建 CompositeProvider
	composite := core.CreateCompositeProvider(localProvider, remoteProvider)

	// 创建使用 CompositeProvider 的 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{
			Name:        "composite",
			Description: "Skill with composite provider",
		},
		Body: `
Composite Skill

Local: <script>local_script</script>
Remote: <script>remote_script</script>
`,
		Provider: composite,
	}

	ctx := context.Background()

	// 执行本地脚本
	result, _ := skill.UseScript(ctx, "local_script", `{}`)
	fmt.Printf("Local: %s\n", result)

	// 执行远程脚本
	result, _ = skill.UseScript(ctx, "remote_script", `{}`)
	fmt.Printf("Remote: %s\n", result)

	// 读取本地参考文档
	ref, _ := skill.ReadReference("local_ref")
	fmt.Printf("Local ref: %s\n", ref)

	// 读取远程参考文档
	ref, _ = skill.ReadReference("remote_ref")
	fmt.Printf("Remote ref: %s\n", ref)

	// 列出所有脚本
	scriptNames, _ := composite.ListScripts(ctx)
	fmt.Printf("All scripts: %v\n", scriptNames)
}

// TestLazyLoadingRemoteProvider 演示懒加载远程 Provider
func TestLazyLoadingRemoteProvider(t *testing.T) {
	loadCount := 0

	// 创建懒加载 Provider
	lazyProvider := core.CreateLazyLoadingProvider(func(ctx context.Context) (resources.ResourceProvider, error) {
		loadCount++
		fmt.Printf("Loading remote provider (count: %d)\n", loadCount)

		// 模拟远程服务
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resources.ScriptCallResponse{Result: `{"lazy":true}`})
		}))

		provider := core.CreateInlineProvider()
		provider.AddScript(resources.NewRemoteScript("lazy_script", resources.NewHTTPRemoteScriptClient(server.URL)))
		return provider, nil
	})

	// 创建使用懒加载 Provider 的 Skill
	skill := &schema.Skill{
		Metadata: &schema.SkillMetadata{
			Name:        "lazy",
			Description: "Skill with lazy loading provider",
		},
		Body:     `Lazy Skill - <script>lazy_script</script>`,
		Provider: lazyProvider,
	}

	ctx := context.Background()

	// 验证还未加载
	if loadCount != 0 {
		t.Errorf("Expected 0 loads, got %d", loadCount)
	}

	// 第一次访问触发加载
	result, _ := skill.UseScript(ctx, "lazy_script", `{}`)
	fmt.Printf("First call: %s\n", result)

	// 验证已加载
	if loadCount != 1 {
		t.Errorf("Expected 1 load, got %d", loadCount)
	}

	// 第二次访问不应该重新加载
	result, _ = skill.UseScript(ctx, "lazy_script", `{}`)
	fmt.Printf("Second call: %s\n", result)

	if loadCount != 1 {
		t.Errorf("Expected still 1 load, got %d", loadCount)
	}
}
