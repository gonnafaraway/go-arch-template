package company

import (
	"context"
	"fmt"

	"go-arch-template/internal/api/domain/company"
)

type MockRepository struct {
	companies map[string]*company.Company
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		companies: make(map[string]*company.Company),
	}
}

func (m *MockRepository) Create(ctx context.Context, c *company.Company) error {
	if c.ID == "" {
		c.ID = fmt.Sprintf("company_%d", len(m.companies)+1)
	}
	m.companies[c.ID] = c
	return nil
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*company.Company, error) {
	c, ok := m.companies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (m *MockRepository) FindAll(ctx context.Context) ([]*company.Company, error) {
	result := make([]*company.Company, 0, len(m.companies))
	for _, c := range m.companies {
		result = append(result, c)
	}
	return result, nil
}

func (m *MockRepository) Update(ctx context.Context, c *company.Company) error {
	if _, ok := m.companies[c.ID]; !ok {
		return ErrNotFound
	}
	m.companies[c.ID] = c
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.companies[id]; !ok {
		return ErrNotFound
	}
	delete(m.companies, id)
	return nil
}


