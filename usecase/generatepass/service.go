package generatepass

import (
	"businesscard-wallet-backend/internal/domain/businesscard"
	"businesscard-wallet-backend/ports"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Service struct {
	templateProvider ports.TemplateProvider
	assetProvider    ports.AssetProvider
	signer           ports.Signer
	packager         ports.Packager
}

func NewService(
	templateProvider ports.TemplateProvider,
	assetProvider ports.AssetProvider,
	signer ports.Signer,
	packager ports.Packager,
) *Service {
	return &Service{
		templateProvider: templateProvider,
		assetProvider:    assetProvider,
		signer:           signer,
		packager:         packager,
	}
}

func (s *Service) Execute(card businesscard.BusinessCard) ([]byte, error) {
	passJSON, err := s.templateProvider.BuildPassJSON(card)
	if err != nil {
		return nil, fmt.Errorf("build pass json: %w", err)
	}

	assets, err := s.assetProvider.LoadAssets(card)
	if err != nil {
		return nil, fmt.Errorf("load assets: %w", err)
	}

	files := make([]ports.PackageFile, 0, 4+len(assets))
	files = append(files, ports.PackageFile{
		Name: "pass.json",
		Data: passJSON,
	})

	for _, asset := range assets {
		files = append(files, ports.PackageFile{
			Name: asset.Name,
			Data: asset.Data,
		})
	}

	manifestBytes, err := buildManifest(files)
	if err != nil {
		return nil, fmt.Errorf("build manifest: %w", err)
	}

	signature, err := s.signer.Sign(manifestBytes)
	if err != nil {
		return nil, fmt.Errorf("sign manifest: %w", err)
	}

	files = append(files,
		ports.PackageFile{
			Name: "manifest.json",
			Data: manifestBytes,
		},
		ports.PackageFile{
			Name: "signature",
			Data: signature,
		},
	)

	pkpass, err := s.packager.Package(files)
	if err != nil {
		return nil, fmt.Errorf("package pass: %w", err)
	}

	return pkpass, nil
}

func buildManifest(files []ports.PackageFile) ([]byte, error) {
	manifest := make(map[string]string, len(files))

	for _, file := range files {
		hash := sha1.Sum(file.Data)
		manifest[file.Name] = hex.EncodeToString(hash[:])
	}

	return json.MarshalIndent(manifest, "", "  ")
}
