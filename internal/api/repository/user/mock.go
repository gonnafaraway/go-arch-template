package user

import (
	"context"
	"fmt"

	"go-arch-template/internal/api/domain/user"
)

type MockRepository struct {
	users map[string]*user.User
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		users: make(map[string]*user.User),
	}
}

func (m *MockRepository) Create(ctx context.Context, u *user.User) error {
	if u.ID == "" {
		u.ID = fmt.Sprintf("user_%d", len(m.users)+1)
	}
	m.users[u.ID] = u
	return nil
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}

func (m *MockRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, ErrNotFound
}

func (m *MockRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, ok := m.users[id]
	return ok, nil
}

func (m *MockRepository) FindAll(ctx context.Context) ([]*user.User, error) {
	result := make([]*user.User, 0, len(m.users))
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, nil
}

func (m *MockRepository) Update(ctx context.Context, u *user.User) error {
	if _, ok := m.users[u.ID]; !ok {
		return ErrNotFound
	}
	m.users[u.ID] = u
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.users[id]; !ok {
		return ErrNotFound
	}
	delete(m.users, id)
	return nil
}



