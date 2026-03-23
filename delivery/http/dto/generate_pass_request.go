package dto

type GeneratePassRequest struct {
	FullName string `json:"fullName"`
	Title    string `json:"title"`
	Company  string `json:"company"`
	Email    string `json:"email"`
	LogoText string `json:"logoText"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
}
