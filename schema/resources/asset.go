package resources

import "fmt"

// todo 暂时没get到它的实际用处

const (
	PNG  AssetExt = "png"
	PPTX AssetExt = "pptx"
	TTF  AssetExt = "ttf"
	PDF  AssetExt = "pdf"
	JPG  AssetExt = "jpg"
	JPEG AssetExt = "jpeg"
	GIF  AssetExt = "gif"
	SVG  AssetExt = "svg"
	WOFF AssetExt = "woff"
	WOFF2 AssetExt = "woff2"
)

type AssetExt string

type Asset struct {
	Name  string   `json:"name"`
	Bytes []byte   `json:"bytes"`
	Ext   AssetExt `json:"ext"`
}

// Size returns the size of the asset in bytes
func (a *Asset) Size() int64 {
	return int64(len(a.Bytes))
}

// String returns a string representation of the asset
func (a *Asset) String() string {
	return fmt.Sprintf("Asset{Name: %s, Ext: %s, Size: %d bytes}", a.Name, a.Ext, a.Size())
}
