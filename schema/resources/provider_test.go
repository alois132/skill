package resources

import (
	"context"
	"testing"
)

func TestInlineProvider(t *testing.T) {
	ctx := context.Background()
	provider := NewInlineProvider()

	// 创建测试脚本
	script := NewEasyScript("test_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"result": "ok"}, nil
	})

	// 创建测试参考文档
	ref := &Reference{Name: "test_ref", Body: "Reference content"}

	// 创建测试资源文件
	asset := &Asset{Name: "test_asset", Bytes: []byte("asset data"), Ext: "txt"}

	// 添加到 provider
	provider.AddScript(script)
	provider.AddReference(ref)
	provider.AddAsset(asset)

	// 测试 GetScript
	gotScript, err := provider.GetScript(ctx, "test_script")
	if err != nil {
		t.Fatalf("Failed to get script: %v", err)
	}
	if gotScript.GetName() != "test_script" {
		t.Errorf("Expected script name 'test_script', got '%s'", gotScript.GetName())
	}

	// 测试 GetReference
	gotRef, err := provider.GetReference(ctx, "test_ref")
	if err != nil {
		t.Fatalf("Failed to get reference: %v", err)
	}
	if gotRef != "Reference content" {
		t.Errorf("Expected reference content 'Reference content', got '%s'", gotRef)
	}

	// 测试 GetAsset
	gotAsset, err := provider.GetAsset(ctx, "test_asset")
	if err != nil {
		t.Fatalf("Failed to get asset: %v", err)
	}
	if string(gotAsset.Bytes) != "asset data" {
		t.Errorf("Expected asset data 'asset data', got '%s'", string(gotAsset.Bytes))
	}

	// 测试 ListScripts
	scriptNames, err := provider.ListScripts(ctx)
	if err != nil {
		t.Fatalf("Failed to list scripts: %v", err)
	}
	if len(scriptNames) != 1 || scriptNames[0] != "test_script" {
		t.Errorf("Expected ['test_script'], got %v", scriptNames)
	}

	// 测试 ListReferences
	refNames, err := provider.ListReferences(ctx)
	if err != nil {
		t.Fatalf("Failed to list references: %v", err)
	}
	if len(refNames) != 1 || refNames[0] != "test_ref" {
		t.Errorf("Expected ['test_ref'], got %v", refNames)
	}

	// 测试 ListAssets
	assetNames, err := provider.ListAssets(ctx)
	if err != nil {
		t.Fatalf("Failed to list assets: %v", err)
	}
	if len(assetNames) != 1 || assetNames[0] != "test_asset" {
		t.Errorf("Expected ['test_asset'], got %v", assetNames)
	}
}

func TestInlineProvider_NotFound(t *testing.T) {
	ctx := context.Background()
	provider := NewInlineProvider()

	// 测试获取不存在的资源
	_, err := provider.GetScript(ctx, "non_existent")
	if err == nil {
		t.Error("Expected error for non-existent script")
	}

	_, err = provider.GetReference(ctx, "non_existent")
	if err == nil {
		t.Error("Expected error for non-existent reference")
	}

	_, err = provider.GetAsset(ctx, "non_existent")
	if err == nil {
		t.Error("Expected error for non-existent asset")
	}
}

func TestCompositeProvider(t *testing.T) {
	ctx := context.Background()

	// 创建两个 provider
	provider1 := NewInlineProvider()
	provider1.AddScript(NewEasyScript("script1", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"from": "provider1"}, nil
	}))
	provider1.AddReference(&Reference{Name: "ref1", Body: "Content from provider1"})

	provider2 := NewInlineProvider()
	provider2.AddScript(NewEasyScript("script2", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"from": "provider2"}, nil
	}))
	provider2.AddReference(&Reference{Name: "ref2", Body: "Content from provider2"})

	// 创建 composite provider
	composite := NewCompositeProvider(provider1, provider2)

	// 测试从第一个 provider 获取
	script1, err := composite.GetScript(ctx, "script1")
	if err != nil {
		t.Fatalf("Failed to get script1: %v", err)
	}
	if script1.GetName() != "script1" {
		t.Errorf("Expected script name 'script1', got '%s'", script1.GetName())
	}

	// 测试从第二个 provider 获取
	script2, err := composite.GetScript(ctx, "script2")
	if err != nil {
		t.Fatalf("Failed to get script2: %v", err)
	}
	if script2.GetName() != "script2" {
		t.Errorf("Expected script name 'script2', got '%s'", script2.GetName())
	}

	// 测试合并列表
	scriptNames, err := composite.ListScripts(ctx)
	if err != nil {
		t.Fatalf("Failed to list scripts: %v", err)
	}
	if len(scriptNames) != 2 {
		t.Errorf("Expected 2 scripts, got %d", len(scriptNames))
	}
}

func TestCompositeProvider_Priority(t *testing.T) {
	ctx := context.Background()

	// 创建两个 provider，都有同名 script
	provider1 := NewInlineProvider()
	provider1.AddScript(NewEasyScript("shared_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"from": "provider1"}, nil
	}))

	provider2 := NewInlineProvider()
	provider2.AddScript(NewEasyScript("shared_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"from": "provider2"}, nil
	}))

	// provider1 优先级更高
	composite := NewCompositeProvider(provider1, provider2)

	script, err := composite.GetScript(ctx, "shared_script")
	if err != nil {
		t.Fatalf("Failed to get script: %v", err)
	}

	// 执行脚本验证来源
	result, err := script.Run(ctx, "{}")
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	// 应该来自 provider1
	if result != `{"from":"provider1"}` {
		t.Errorf("Expected result from provider1, got '%s'", result)
	}
}

func TestCachingProvider(t *testing.T) {
	ctx := context.Background()

	// 创建基础 provider
	baseProvider := NewInlineProvider()
	callCount := 0
	baseProvider.AddScript(NewEasyScript("test_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		callCount++
		return map[string]interface{}{"call": callCount}, nil
	}))
	baseProvider.AddReference(&Reference{Name: "test_ref", Body: "Content"})

	// 创建 caching provider
	caching := NewCachingProvider(baseProvider)

	// 第一次获取脚本
	script1, err := caching.GetScript(ctx, "test_script")
	if err != nil {
		t.Fatalf("Failed to get script: %v", err)
	}

	// 第二次获取脚本（应该从缓存）
	script2, err := caching.GetScript(ctx, "test_script")
	if err != nil {
		t.Fatalf("Failed to get script: %v", err)
	}

	// 验证是同一个对象
	if script1 != script2 {
		t.Error("Expected cached script to be the same object")
	}

	// 测试缓存参考文档
	ref1, err := caching.GetReference(ctx, "test_ref")
	if err != nil {
		t.Fatalf("Failed to get reference: %v", err)
	}

	ref2, err := caching.GetReference(ctx, "test_ref")
	if err != nil {
		t.Fatalf("Failed to get reference: %v", err)
	}

	if ref1 != ref2 {
		t.Error("Expected cached reference to be the same")
	}

	// 清除缓存后应该重新获取
	caching.ClearCache()
	_, err = caching.GetScript(ctx, "test_script")
	if err != nil {
		t.Fatalf("Failed to get script after clear: %v", err)
	}
}

func TestLazyLoadingProvider(t *testing.T) {
	ctx := context.Background()

	loadCount := 0
	loader := func(ctx context.Context) (ResourceProvider, error) {
		loadCount++
		provider := NewInlineProvider()
		provider.AddScript(NewEasyScript("lazy_script", func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{"loaded": true}, nil
		}))
		return provider, nil
	}

	// 创建懒加载 provider
	lazy := NewLazyLoadingProvider(loader)

	// 验证还未加载
	if loadCount != 0 {
		t.Errorf("Expected 0 loads, got %d", loadCount)
	}

	// 第一次访问触发加载
	script, err := lazy.GetScript(ctx, "lazy_script")
	if err != nil {
		t.Fatalf("Failed to get script: %v", err)
	}

	// 验证已加载
	if loadCount != 1 {
		t.Errorf("Expected 1 load, got %d", loadCount)
	}

	if script.GetName() != "lazy_script" {
		t.Errorf("Expected script name 'lazy_script', got '%s'", script.GetName())
	}

	// 第二次访问不应该重新加载
	_, err = lazy.GetScript(ctx, "lazy_script")
	if err != nil {
		t.Fatalf("Failed to get script: %v", err)
	}

	if loadCount != 1 {
		t.Errorf("Expected still 1 load, got %d", loadCount)
	}
}
