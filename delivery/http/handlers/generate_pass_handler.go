package handlers

import (
	"businesscard-wallet-backend/delivery/http/dto"
	"businesscard-wallet-backend/internal/domain/businesscard"
	"businesscard-wallet-backend/usecase/generatepass"
	"encoding/json"
	"errors"
	"net/http"
)

type GeneratePassHandler struct {
	service *generatepass.Service
}

func NewGeneratePassHandler(service *generatepass.Service) *GeneratePassHandler {
	return &GeneratePassHandler{service: service}
}

func (h *GeneratePassHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.GeneratePassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	card, err := businesscard.New(
		req.FullName,
		req.Title,
		req.Company,
		req.LogoText,
		req.Email,
		req.Phone,
		req.Website,
	)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	pkpass, err := h.service.Execute(card)
	if err != nil {
		http.Error(w, "failed to generate pass", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.apple.pkpass")
	w.Header().Set("Content-Disposition", `attachment; filename="business-card.pkpass"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pkpass)
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, businesscard.ErrFullNameRequired),
		errors.Is(err, businesscard.ErrTitleRequired),
		errors.Is(err, businesscard.ErrCompanyRequired),
		errors.Is(err, businesscard.ErrContactRequired):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, "invalid business card data", http.StatusBadRequest)
	}
}
