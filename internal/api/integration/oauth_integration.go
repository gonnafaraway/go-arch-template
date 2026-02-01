package integration

import (
	"context"
)

type oauthIntegration struct {
	// In a real application, there would be an OAuth client here
}

func NewOAuthIntegration() OAuthIntegration {
	return &oauthIntegration{}
}

func (i *oauthIntegration) ValidateToken(ctx context.Context, token string) (bool, error) {
	// Mock token validation
	if token == "" {
		return false, nil
	}
	return true, nil
}

func (i *oauthIntegration) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	// Mock user info retrieval
	return &UserInfo{
		UserID:  "user_1",
		Email:   "user@example.com",
		Name:    "Test User",
		Company: "Test Company",
	}, nil
}
