package ports

import "businesscard-wallet-backend/internal/domain/businesscard"

type TemplateProvider interface {
	BuildPassJSON(card businesscard.BusinessCard) ([]byte, error)
}
