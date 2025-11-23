package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-arch-template/internal/api/usecase/company"
)

// HandleCompanies handles company list and creation
// @Summary List or create companies
// @Description Get list of all companies or create a new company
// @Tags companies
// @Accept json
// @Produce json
// @Param company body company.CreateCompanyRequest false "Company data"
// @Success 200 {array} company.CompanyResponse "List of companies"
// @Success 201 {object} company.CompanyResponse "Created company"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/companies [get]
// @Router /api/companies [post]

type CompanyHandler struct {
	useCase *company.CompanyUseCase
}

func NewCompanyHandler(useCase *company.CompanyUseCase) *CompanyHandler {
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
		var req company.CreateCompanyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request body", err)
			return
		}
		companyResp, err := h.useCase.CreateCompany(ctx, req)
		if err != nil {
			// Проверяем тип ошибки для правильного статус кода
			if isValidationError(err) {
				respondError(w, http.StatusBadRequest, "validation failed", err)
			} else {
				respondError(w, http.StatusInternalServerError, "failed to create company", err)
			}
			return
		}
		respondJSON(w, http.StatusCreated, companyResp)

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
// @Success 200 {object} company.CompanyResponse "Company information"
// @Failure 404 {object} map[string]string "Company not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/companies/{id} [get]
func (h *CompanyHandler) HandleCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := strings.TrimPrefix(r.URL.Path, "/api/companies/")

	switch r.Method {
	case http.MethodGet:
		companyResp, err := h.useCase.GetCompany(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		respondJSON(w, http.StatusOK, companyResp)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
