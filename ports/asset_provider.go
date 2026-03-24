package ports

import "businesscard-wallet-backend/internal/domain/businesscard"

type PassAsset struct {
	Name string
	Data []byte
}

type AssetProvider interface {
	LoadAssets(card businesscard.BusinessCard) ([]PassAsset, error)
}
