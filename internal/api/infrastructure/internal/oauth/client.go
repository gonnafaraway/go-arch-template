package oauth

import (
	"context"
)

type Client interface {
	ValidateToken(ctx context.Context, token string) (*TokenInfo, error)
	GetUserInfo(ctx context.Context, token string) (*UserInfo, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
}

type TokenInfo struct {
	Valid   bool   `json:"valid"`
	UserID  string `json:"user_id"`
	Expires int64  `json:"expires"`
}

type UserInfo struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Company string `json:"company"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type MockClient struct {
	tokens map[string]*TokenInfo
}

func NewMockClient() *MockClient {
	return &MockClient{
		tokens: make(map[string]*TokenInfo),
	}
}

func (m *MockClient) ValidateToken(ctx context.Context, token string) (*TokenInfo, error) {
	info, ok := m.tokens[token]
	if !ok {
		return &TokenInfo{Valid: false}, nil
	}
	return info, nil
}

func (m *MockClient) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	info, ok := m.tokens[token]
	if !ok || !info.Valid {
		return nil, ErrInvalidToken
	}
	return &UserInfo{
		UserID:  info.UserID,
		Email:   "user@example.com",
		Name:    "Test User",
		Company: "Test Company",
	}, nil
}

func (m *MockClient) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	return &TokenResponse{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
		ExpiresIn:    3600,
	}, nil
}
