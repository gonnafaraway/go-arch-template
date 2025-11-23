package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-arch-template/internal/api/usecase"
)

type CompanyHandler struct {
	useCase *usecase.CompanyUseCase
}

func NewCompanyHandler(useCase *usecase.CompanyUseCase) *CompanyHandler {
	return &CompanyHandler{
		useCase: useCase,
	}
}

func (h *CompanyHandler) HandleCompanies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		companies, err := h.useCase.ListCompanies(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusOK, companies)

	case http.MethodPost:
		var req usecase.CreateCompanyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		company, err := h.useCase.CreateCompany(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusCreated, company)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CompanyHandler) HandleCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := strings.TrimPrefix(r.URL.Path, "/api/companies/")

	switch r.Method {
	case http.MethodGet:
		company, err := h.useCase.GetCompany(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		respondJSON(w, http.StatusOK, company)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

