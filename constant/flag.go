package constant

const (
	// 自包含格式（推荐用于 Skill Body）
	// 格式：<script>name</script>
	ScriptSelfContained    = "<script>%s</script>"
	ReferenceSelfContained = "<reference>%s</reference>"
	AssetSelfContained     = "<asset>%s</asset>"

	// 带 src 属性的完整格式（用于生成文档或特殊场景）
	ScriptFlag      = `<script src="%s">%s</script>`
	ScriptUsageFlag = `<script_usage src="%s">%s</script_usage>`
	ReferenceFlag   = `<reference src="%s">%s</reference>`
	AssetFlag       = `<asset src="%s" type="%s">%s</asset>`
)
