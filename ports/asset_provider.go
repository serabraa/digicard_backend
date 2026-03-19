package ports

type PassAsset struct {
	Name string
	Data []byte
}

type AssetProvider interface {
	LoadAssets() ([]PassAsset, error)
}
