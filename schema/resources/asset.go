package resources

// todo 暂时没get到它的实际用处

const (
	PNG  AssetExt = "png"
	PPTX AssetExt = "pptx"
	TTF  AssetExt = "ttf"
)

type AssetExt string

type Asset struct {
	Name  string   `json:"name"`
	Bytes []byte   `json:"bytes"`
	Ext   AssetExt `json:"ext"`
}
