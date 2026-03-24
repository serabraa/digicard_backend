package dto

type GeneratePassRequest struct {
	FullName        string `json:"fullName"`
	Title           string `json:"title"`
	Company         string `json:"company"`
	LogoText        string `json:"logoText"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Website         string `json:"website"`
	QRCodeContent   string `json:"qrCodeContent"`
	BackgroundColor string `json:"backgroundColor"`
	ForegroundColor string `json:"foregroundColor"`
	LabelColor      string `json:"labelColor"`
}
