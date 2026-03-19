package dto

type GeneratePassRequest struct {
	FullName string `json:"fullName"`
	Title    string `json:"title"`
	Company  string `json:"company"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
}
