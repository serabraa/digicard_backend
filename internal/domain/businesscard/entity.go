package businesscard

import "strings"

type BusinessCard struct {
	FullName             string
	Title                string
	Company              string
	LogoText             string
	Email                string
	Phone                string
	Website              string
	QRCodeContent        string
	BackgroundColor      string
	ForegroundColor      string
	LabelColor           string
	LogoImageBase64      string
	ThumbnailImageBase64 string
}

func New(
	fullName,
	title,
	company,
	logoText,
	email,
	phone,
	website,
	qrCodeContent,
	backgroundColor,
	foregroundColor,
	labelColor,
	logoImageBase64,
	thumbnailImageBase64 string,
) (BusinessCard, error) {
	card := BusinessCard{
		FullName:             strings.TrimSpace(fullName),
		Title:                strings.TrimSpace(title),
		Company:              strings.TrimSpace(company),
		LogoText:             strings.TrimSpace(logoText),
		Email:                strings.TrimSpace(email),
		Phone:                strings.TrimSpace(phone),
		Website:              strings.TrimSpace(website),
		QRCodeContent:        strings.TrimSpace(qrCodeContent),
		BackgroundColor:      strings.TrimSpace(backgroundColor),
		ForegroundColor:      strings.TrimSpace(foregroundColor),
		LabelColor:           strings.TrimSpace(labelColor),
		LogoImageBase64:      strings.TrimSpace(logoImageBase64),
		ThumbnailImageBase64: strings.TrimSpace(thumbnailImageBase64),
	}

	if err := card.Validate(); err != nil {
		return BusinessCard{}, err
	}

	return card, nil
}

func (c BusinessCard) Validate() error {
	if c.FullName == "" {
		return ErrFullNameRequired
	}
	return nil
}
