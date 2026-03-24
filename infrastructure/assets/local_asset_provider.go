package assets

import (
	"businesscard-wallet-backend/internal/domain/businesscard"
	"businesscard-wallet-backend/ports"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type LocalAssetProvider struct {
	dir string
}

func NewLocalAssetProvider(dir string) *LocalAssetProvider {
	return &LocalAssetProvider{dir: dir}
}

func (p *LocalAssetProvider) LoadAssets(card businesscard.BusinessCard) ([]ports.PassAsset, error) {
	assets := make([]ports.PassAsset, 0, 3)

	iconData, err := os.ReadFile(filepath.Join(p.dir, "icon.png"))
	if err != nil {
		return nil, fmt.Errorf("read asset icon.png: %w", err)
	}
	assets = append(assets, ports.PassAsset{
		Name: "icon.png",
		Data: iconData,
	})

	logoData, err := p.resolveImage(card.LogoImageBase64, "logo.png")
	if err != nil {
		return nil, fmt.Errorf("resolve logo.png: %w", err)
	}
	if len(logoData) > 0 {
		assets = append(assets, ports.PassAsset{
			Name: "logo.png",
			Data: logoData,
		})
	}

	thumbnailData, err := p.resolveImage(card.ThumbnailImageBase64, "thumbnail.png")
	if err != nil {
		return nil, fmt.Errorf("resolve thumbnail.png: %w", err)
	}
	if len(thumbnailData) > 0 {
		assets = append(assets, ports.PassAsset{
			Name: "thumbnail.png",
			Data: thumbnailData,
		})
	}

	return assets, nil
}

func (p *LocalAssetProvider) resolveImage(base64Value string, fallbackFileName string) ([]byte, error) {
	if strings.TrimSpace(base64Value) != "" {
		data, err := decodeBase64Image(base64Value)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	fullPath := filepath.Join(p.dir, fallbackFileName)
	data, err := os.ReadFile(fullPath)
	if err == nil {
		return data, nil
	}
	if os.IsNotExist(err) {
		return nil, nil
	}
	return nil, err
}

func decodeBase64Image(value string) ([]byte, error) {
	clean := strings.TrimSpace(value)

	if commaIndex := strings.Index(clean, ","); commaIndex != -1 {
		clean = clean[commaIndex+1:]
	}

	data, err := base64.StdEncoding.DecodeString(clean)
	if err == nil {
		return data, nil
	}

	data, err = base64.RawStdEncoding.DecodeString(clean)
	if err == nil {
		return data, nil
	}

	return nil, fmt.Errorf("invalid base64 image data")
}
