package template

import (
	"businesscard-wallet-backend/internal/domain/businesscard"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
)

type StaticTemplateProvider struct {
	organizationName string
	passTypeID       string
	teamID           string
	backgroundColor  string
	labelColor       string
	foregroundColor  string
}

func NewStaticTemplateProvider(
	organizationName string,
	passTypeID string,
	teamID string,
	backgroundColor string,
	labelColor string,
	foregroundColor string,
) *StaticTemplateProvider {
	return &StaticTemplateProvider{
		organizationName: organizationName,
		passTypeID:       passTypeID,
		teamID:           teamID,
		backgroundColor:  backgroundColor,
		labelColor:       labelColor,
		foregroundColor:  foregroundColor,
	}
}

func (p *StaticTemplateProvider) BuildPassJSON(card businesscard.BusinessCard) ([]byte, error) {
	payload := map[string]any{
		"formatVersion":      1,
		"passTypeIdentifier": p.passTypeID,
		"serialNumber":       generateSerial(card),
		"teamIdentifier":     p.teamID,
		"organizationName":   p.organizationName,
		"description":        "Digital business card",
		"backgroundColor":    p.backgroundColor,
		"labelColor":         p.labelColor,
		"foregroundColor":    p.foregroundColor,
		"generic": map[string]any{
			"primaryFields": []map[string]any{
				{
					"key":   "name",
					"label": "Name",
					"value": card.FullName,
				},
			},
			"secondaryFields": []map[string]any{
				{
					"key":   "title",
					"label": "Title",
					"value": card.Title,
				},
				{
					"key":   "company",
					"label": "Company",
					"value": card.Company,
				},
			},
			"auxiliaryFields": buildAuxiliaryFields(card),
		},
	}

	return json.MarshalIndent(payload, "", "  ")
}

func buildAuxiliaryFields(card businesscard.BusinessCard) []map[string]any {
	fields := make([]map[string]any, 0, 3)

	if card.Email != "" {
		fields = append(fields, map[string]any{
			"key":   "email",
			"label": "Email",
			"value": card.Email,
		})
	}

	if card.Phone != "" {
		fields = append(fields, map[string]any{
			"key":   "phone",
			"label": "Phone",
			"value": card.Phone,
		})
	}

	if card.Website != "" {
		fields = append(fields, map[string]any{
			"key":   "website",
			"label": "Website",
			"value": card.Website,
		})
	}

	return fields
}

func generateSerial(card businesscard.BusinessCard) string {
	raw := card.FullName + "|" + card.Email + "|" + card.Phone + "|" + card.Company

	hash := sha1.Sum([]byte(raw))
	return hex.EncodeToString(hash[:])
}
