package util

import (
	"regexp"
	"strings"
)

// XMLTag 表示解析出的 XML 标记
type XMLTag struct {
	TagName string // 标记名：script, reference, asset
	Content string // 标记内容（如 "init_skill", "usage_guide"）
}

// ParseXMLTags 从文本中解析所有 XML 标记
// 支持格式：<script>name</script> 或 <reference>name</reference>
func ParseXMLTags(body string) []XMLTag {
	if body == "" {
		return nil
	}

	// 正则匹配 XML 标记：支持 script, reference, asset
	// 格式：<tag>content</tag>
	pattern := `<(script|reference|asset)>([^<]+)</(script|reference|asset)>`
	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatch(body, -1)
	if matches == nil {
		return nil
	}

	tags := make([]XMLTag, 0, len(matches))
	for _, match := range matches {
		if len(match) >= 3 {
			tag := XMLTag{
				TagName: match[1],
				Content: strings.TrimSpace(match[2]),
			}
			tags = append(tags, tag)
		}
	}

	return tags
}

// ExtractScriptNames 从 Body 中提取所有脚本名称
func ExtractScriptNames(body string) []string {
	tags := ParseXMLTags(body)
	if tags == nil {
		return nil
	}

	scriptNames := []string{}
	for _, tag := range tags {
		if tag.TagName == "script" {
			scriptNames = append(scriptNames, tag.Content)
		}
	}

	if len(scriptNames) == 0 {
		return nil
	}
	return scriptNames
}

// ExtractReferenceNames 从 Body 中提取所有参考文献名称
func ExtractReferenceNames(body string) []string {
	tags := ParseXMLTags(body)
	if tags == nil {
		return nil
	}

	refNames := []string{}
	for _, tag := range tags {
		if tag.TagName == "reference" {
			refNames = append(refNames, tag.Content)
		}
	}

	if len(refNames) == 0 {
		return nil
	}
	return refNames
}

// ExtractAssetNames 从 Body 中提取所有资产名称
func ExtractAssetNames(body string) []string {
	tags := ParseXMLTags(body)
	if tags == nil {
		return nil
	}

	assetNames := []string{}
	for _, tag := range tags {
		if tag.TagName == "asset" {
			assetNames = append(assetNames, tag.Content)
		}
	}

	if len(assetNames) == 0 {
		return nil
	}
	return assetNames
}

// HasXMLTags 检查 Body 中是否包含 XML 标记
func HasXMLTags(body string) bool {
	return len(ParseXMLTags(body)) > 0
}
