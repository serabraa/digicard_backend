package businesscard

import "strings"

type BusinessCard struct {
	FullName string
	Title    string
	Company  string
	LogoText string
	Email    string
	Phone    string
	Website  string
}

func New(fullName, title, company, logoText, email, phone, website string) (BusinessCard, error) {
	card := BusinessCard{
		FullName: strings.TrimSpace(fullName),
		Title:    strings.TrimSpace(title),
		Company:  strings.TrimSpace(company),
		LogoText: strings.TrimSpace(logoText),
		Email:    strings.TrimSpace(email),
		Phone:    strings.TrimSpace(phone),
		Website:  strings.TrimSpace(website),
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
	if c.Title == "" {
		return ErrTitleRequired
	}
	if c.Company == "" {
		return ErrCompanyRequired
	}
	if c.Email == "" && c.Phone == "" && c.Website == "" {
		return ErrContactRequired
	}

	return nil
}
