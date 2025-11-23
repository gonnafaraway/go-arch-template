package usersservice

import (
	"context"
)

type Client interface {
	GetUser(ctx context.Context, userID string) (*UserResponse, error)
	ValidateUser(ctx context.Context, userID string) (bool, error)
	SyncUser(ctx context.Context, userID string) error
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CompanyID string `json:"company_id"`
	Active    bool   `json:"active"`
}

type MockClient struct {
	users map[string]*UserResponse
}

func NewMockClient() *MockClient {
	return &MockClient{
		users: make(map[string]*UserResponse),
	}
}

func (m *MockClient) GetUser(ctx context.Context, userID string) (*UserResponse, error) {
	user, ok := m.users[userID]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (m *MockClient) ValidateUser(ctx context.Context, userID string) (bool, error) {
	_, ok := m.users[userID]
	return ok, nil
}

func (m *MockClient) SyncUser(ctx context.Context, userID string) error {
	// Мок синхронизации
	return nil
}
