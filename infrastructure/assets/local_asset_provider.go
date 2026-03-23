package assets

import (
	"businesscard-wallet-backend/ports"
	"fmt"
	"os"
	"path/filepath"
)

type LocalAssetProvider struct {
	dir string
}

func NewLocalAssetProvider(dir string) *LocalAssetProvider {
	return &LocalAssetProvider{dir: dir}
}

func (p *LocalAssetProvider) LoadAssets() ([]ports.PassAsset, error) {
	requiredFiles := []string{"icon.png", "logo.png"}
	optionalFiles := []string{"thumbnail.png"}

	assets := make([]ports.PassAsset, 0, len(requiredFiles)+len(optionalFiles))

	for _, name := range requiredFiles {
		fullPath := filepath.Join(p.dir, name)

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("read asset %s: %w", name, err)
		}

		assets = append(assets, ports.PassAsset{
			Name: name,
			Data: data,
		})
	}

	for _, name := range optionalFiles {
		fullPath := filepath.Join(p.dir, name)

		data, err := os.ReadFile(fullPath)
		if err == nil {
			assets = append(assets, ports.PassAsset{
				Name: name,
				Data: data,
			})
		}
	}

	return assets, nil
}
