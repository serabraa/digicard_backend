package businesscard

import "errors"

var (
	ErrFullNameRequired = errors.New("full name is required")
	ErrTitleRequired    = errors.New("title is required")
	ErrCompanyRequired  = errors.New("company is required")
	ErrContactRequired  = errors.New("at least one contact method is required")
)
