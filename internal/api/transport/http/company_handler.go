package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-arch-template/internal/api/usecase"
)

// HandleCompanies handles company list and creation
// @Summary List or create companies
// @Description Get list of all companies or create a new company
// @Tags companies
// @Accept json
// @Produce json
// @Param company body usecase.CreateCompanyRequest false "Company data"
// @Success 200 {array} usecase.CompanyResponse "List of companies"
// @Success 201 {object} usecase.CompanyResponse "Created company"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/companies [get]
// @Router /api/companies [post]

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

// HandleCompany handles single company operations
// @Summary Get company by ID
// @Description Get company information by ID
// @Tags companies
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} usecase.CompanyResponse "Company information"
// @Failure 404 {object} map[string]string "Company not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/companies/{id} [get]
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

