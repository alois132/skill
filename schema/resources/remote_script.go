package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// RemoteScriptClient 远程脚本客户端接口
// 定义了调用远程脚本所需的方法
type RemoteScriptClient interface {
	// Call 调用远程脚本
	// scriptName: 脚本名称
	// args: JSON 格式的参数
	// 返回: JSON 格式的结果
	Call(ctx context.Context, scriptName string, args string) (string, error)
}

// RemoteScript 远程脚本实现
// 通过 RemoteScriptClient 调用远程服务执行脚本
type RemoteScript struct {
	Name   string
	Usage  string
	Client RemoteScriptClient
}

// Run 执行远程脚本
func (s *RemoteScript) Run(ctx context.Context, args string) (result string, err error) {
	if s.Client == nil {
		return "", errors.New("remote script client not configured")
	}
	return s.Client.Call(ctx, s.Name, args)
}

// GetName 获取脚本名称
func (s *RemoteScript) GetName() string {
	return s.Name
}

// GetUsage 获取脚本使用说明
func (s *RemoteScript) GetUsage() string {
	return s.Usage
}

// NewRemoteScript 创建一个新的远程脚本
func NewRemoteScript(name string, client RemoteScriptClient) *RemoteScript {
	return &RemoteScript{
		Name:   name,
		Client: client,
		Usage:  fmt.Sprintf("Remote script: %s", name),
	}
}

// WithUsage 设置脚本使用说明
func (s *RemoteScript) WithUsage(usage string) *RemoteScript {
	s.Usage = usage
	return s
}

// Ensure RemoteScript implements Script
var _ Script = (*RemoteScript)(nil)

// HTTPRemoteScriptClient 基于 HTTP 的远程脚本客户端
type HTTPRemoteScriptClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Headers    map[string]string
}

// HTTPClientOption HTTP 客户端配置选项
type HTTPClientOption func(*HTTPRemoteScriptClient)

// NewHTTPRemoteScriptClient 创建一个新的 HTTP 远程脚本客户端
// baseURL: 远程服务的基础 URL，例如 "http://localhost:8080/api/scripts"
func NewHTTPRemoteScriptClient(baseURL string, opts ...HTTPClientOption) *HTTPRemoteScriptClient {
	client := &HTTPRemoteScriptClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Headers: make(map[string]string),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// WithTimeout 设置 HTTP 请求超时
func WithTimeout(timeout time.Duration) HTTPClientOption {
	return func(c *HTTPRemoteScriptClient) {
		c.HTTPClient.Timeout = timeout
	}
}

// WithHeader 添加自定义请求头
func WithHeader(key, value string) HTTPClientOption {
	return func(c *HTTPRemoteScriptClient) {
		c.Headers[key] = value
	}
}

// WithHTTPClient 设置自定义 HTTP 客户端
func WithHTTPClient(httpClient *http.Client) HTTPClientOption {
	return func(c *HTTPRemoteScriptClient) {
		c.HTTPClient = httpClient
	}
}

// ScriptCallRequest HTTP 脚本调用请求
type ScriptCallRequest struct {
	ScriptName string `json:"script_name"`
	Args       string `json:"args"`
}

// ScriptCallResponse HTTP 脚本调用响应
type ScriptCallResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// Call 通过 HTTP 调用远程脚本
func (c *HTTPRemoteScriptClient) Call(ctx context.Context, scriptName string, args string) (string, error) {
	reqBody := ScriptCallRequest{
		ScriptName: scriptName,
		Args:       args,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/%s", c.BaseURL, scriptName)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("remote script returned error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// 尝试解析为 ScriptCallResponse
	var callResp ScriptCallResponse
	if err := json.Unmarshal(body, &callResp); err == nil && callResp.Error != "" {
		return "", errors.New(callResp.Error)
	}

	// 如果解析失败或没有 error 字段，直接返回 body 作为结果
	if callResp.Result != "" {
		return callResp.Result, nil
	}

	return string(body), nil
}

// Ensure HTTPRemoteScriptClient implements RemoteScriptClient
var _ RemoteScriptClient = (*HTTPRemoteScriptClient)(nil)

// MockRemoteScriptClient 用于测试的模拟远程脚本客户端
type MockRemoteScriptClient struct {
	Handlers map[string]func(ctx context.Context, args string) (string, error)
}

// NewMockRemoteScriptClient 创建一个新的模拟远程脚本客户端
func NewMockRemoteScriptClient() *MockRemoteScriptClient {
	return &MockRemoteScriptClient{
		Handlers: make(map[string]func(ctx context.Context, args string) (string, error)),
	}
}

// Register 注册一个脚本处理器
func (m *MockRemoteScriptClient) Register(scriptName string, handler func(ctx context.Context, args string) (string, error)) {
	m.Handlers[scriptName] = handler
}

// Call 调用模拟的远程脚本
func (m *MockRemoteScriptClient) Call(ctx context.Context, scriptName string, args string) (string, error) {
	handler, ok := m.Handlers[scriptName]
	if !ok {
		return "", fmt.Errorf("script not found: %s", scriptName)
	}
	return handler(ctx, args)
}

// Ensure MockRemoteScriptClient implements RemoteScriptClient
var _ RemoteScriptClient = (*MockRemoteScriptClient)(nil)
