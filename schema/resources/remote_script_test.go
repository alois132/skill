package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRemoteScript(t *testing.T) {
	// 创建 mock 客户端
	mockClient := NewMockRemoteScriptClient()
	mockClient.Register("test_script", func(ctx context.Context, args string) (string, error) {
		return `{"result":"success"}`, nil
	})

	// 创建远程脚本
	script := NewRemoteScript("test_script", mockClient)

	// 测试 GetName
	if script.GetName() != "test_script" {
		t.Errorf("Expected name 'test_script', got '%s'", script.GetName())
	}

	// 测试 Run
	ctx := context.Background()
	result, err := script.Run(ctx, `{}`)
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}
	if result != `{"result":"success"}` {
		t.Errorf("Expected result '{\"result\":\"success\"}', got '%s'", result)
	}
}

func TestRemoteScript_NoClient(t *testing.T) {
	script := NewRemoteScript("test_script", nil)

	ctx := context.Background()
	_, err := script.Run(ctx, `{}`)
	if err == nil {
		t.Error("Expected error when client is nil")
	}
}

func TestRemoteScript_WithUsage(t *testing.T) {
	mockClient := NewMockRemoteScriptClient()
	script := NewRemoteScript("test_script", mockClient).WithUsage("Custom usage")

	if script.GetUsage() != "Custom usage" {
		t.Errorf("Expected usage 'Custom usage', got '%s'", script.GetUsage())
	}
}

func TestMockRemoteScriptClient(t *testing.T) {
	client := NewMockRemoteScriptClient()

	// 注册脚本
	client.Register("script1", func(ctx context.Context, args string) (string, error) {
		return "result1", nil
	})
	client.Register("script2", func(ctx context.Context, args string) (string, error) {
		return "result2", nil
	})

	ctx := context.Background()

	// 调用 script1
	result, err := client.Call(ctx, "script1", `{}`)
	if err != nil {
		t.Fatalf("Failed to call script1: %v", err)
	}
	if result != "result1" {
		t.Errorf("Expected 'result1', got '%s'", result)
	}

	// 调用 script2
	result, err = client.Call(ctx, "script2", `{}`)
	if err != nil {
		t.Fatalf("Failed to call script2: %v", err)
	}
	if result != "result2" {
		t.Errorf("Expected 'result2', got '%s'", result)
	}

	// 调用不存在的脚本
	_, err = client.Call(ctx, "non_existent", `{}`)
	if err == nil {
		t.Error("Expected error for non-existent script")
	}
}

func TestHTTPRemoteScriptClient(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// 验证 Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		// 解析请求
		var req ScriptCallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		// 返回响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if req.ScriptName == "test_script" {
			resp := ScriptCallResponse{Result: `{"success":true}`}
			json.NewEncoder(w).Encode(resp)
		} else {
			resp := ScriptCallResponse{Error: "script not found"}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	// 创建 HTTP 客户端
	client := NewHTTPRemoteScriptClient(server.URL)

	ctx := context.Background()

	// 测试成功调用
	result, err := client.Call(ctx, "test_script", `{"key":"value"}`)
	if err != nil {
		t.Fatalf("Failed to call script: %v", err)
	}
	if result != `{"success":true}` {
		t.Errorf("Expected '{\"success\":true}', got '%s'", result)
	}

	// 测试错误响应
	_, err = client.Call(ctx, "non_existent", `{}`)
	if err == nil {
		t.Error("Expected error for non-existent script")
	}
}

func TestHTTPRemoteScriptClient_WithOptions(t *testing.T) {
	// 创建测试服务器，验证自定义 header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证自定义 header
		if r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Errorf("Expected X-Custom-Header 'custom-value', got '%s'", r.Header.Get("X-Custom-Header"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ScriptCallResponse{Result: "ok"})
	}))
	defer server.Close()

	// 创建带自定义选项的 HTTP 客户端
	client := NewHTTPRemoteScriptClient(
		server.URL,
		WithHeader("X-Custom-Header", "custom-value"),
		WithTimeout(10*time.Second),
	)

	ctx := context.Background()
	_, err := client.Call(ctx, "test", `{}`)
	if err != nil {
		t.Fatalf("Failed to call script: %v", err)
	}
}

func TestHTTPRemoteScriptClient_ServerError(t *testing.T) {
	// 创建返回 500 错误的测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewHTTPRemoteScriptClient(server.URL)

	ctx := context.Background()
	_, err := client.Call(ctx, "test", `{}`)
	if err == nil {
		t.Error("Expected error for server error")
	}
}

func TestHTTPRemoteScriptClient_InvalidJSON(t *testing.T) {
	// 创建返回无效 JSON 的测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not valid json"))
	}))
	defer server.Close()

	client := NewHTTPRemoteScriptClient(server.URL)

	ctx := context.Background()
	// 无效 JSON 应该作为原始字符串返回
	result, err := client.Call(ctx, "test", `{}`)
	if err != nil {
		t.Fatalf("Failed to call script: %v", err)
	}
	if result != "not valid json" {
		t.Errorf("Expected 'not valid json', got '%s'", result)
	}
}

func TestHTTPRemoteScriptClient_ContextTimeout(t *testing.T) {
	// 创建慢速服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 创建带超时的客户端
	client := NewHTTPRemoteScriptClient(server.URL, WithTimeout(10*time.Millisecond))

	ctx := context.Background()
	_, err := client.Call(ctx, "test", `{}`)
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestHTTPRemoteScriptClient_URLConstruction(t *testing.T) {
	var capturedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.Path
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`"ok"`))
	}))
	defer server.Close()

	client := NewHTTPRemoteScriptClient(server.URL)

	ctx := context.Background()
	client.Call(ctx, "my_script", `{}`)

	expectedURL := "/my_script"
	if capturedURL != expectedURL {
		t.Errorf("Expected URL path '%s', got '%s'", expectedURL, capturedURL)
	}
}

// TestRemoteScriptIntegration 集成测试：完整的远程脚本调用流程
func TestRemoteScriptIntegration(t *testing.T) {
	// 创建模拟的远程服务
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ScriptCallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// 模拟不同的脚本行为
		switch req.ScriptName {
		case "add":
			var args struct {
				A int `json:"a"`
				B int `json:"b"`
			}
			json.Unmarshal([]byte(req.Args), &args)
			result := fmt.Sprintf(`{"result":%d}`, args.A+args.B)
			json.NewEncoder(w).Encode(ScriptCallResponse{Result: result})

		case "echo":
			json.NewEncoder(w).Encode(ScriptCallResponse{Result: req.Args})

		default:
			json.NewEncoder(w).Encode(ScriptCallResponse{Error: "unknown script"})
		}
	}))
	defer server.Close()

	// 创建 HTTP 客户端和远程脚本
	httpClient := NewHTTPRemoteScriptClient(server.URL)
	addScript := NewRemoteScript("add", httpClient).WithUsage("Add two numbers")
	echoScript := NewRemoteScript("echo", httpClient).WithUsage("Echo the input")

	ctx := context.Background()

	// 测试 add 脚本
	result, err := addScript.Run(ctx, `{"a":5,"b":3}`)
	if err != nil {
		t.Fatalf("Failed to run add script: %v", err)
	}
	if result != `{"result":8}` {
		t.Errorf("Expected '{\"result\":8}', got '%s'", result)
	}

	// 测试 echo 脚本
	result, err = echoScript.Run(ctx, `{"message":"hello"}`)
	if err != nil {
		t.Fatalf("Failed to run echo script: %v", err)
	}
	if result != `{"message":"hello"}` {
		t.Errorf("Expected '{\"message\":\"hello\"}', got '%s'", result)
	}
}
