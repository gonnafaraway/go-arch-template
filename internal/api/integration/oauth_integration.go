package integration

import (
	"context"
)

type oauthIntegration struct {
	// В реальном приложении здесь будет клиент OAuth
}

func NewOAuthIntegration() OAuthIntegration {
	return &oauthIntegration{}
}

func (i *oauthIntegration) ValidateToken(ctx context.Context, token string) (bool, error) {
	// Мок валидации токена
	if token == "" {
		return false, nil
	}
	return true, nil
}

func (i *oauthIntegration) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	// Мок получения информации о пользователе
	return &UserInfo{
		UserID:  "user_1",
		Email:   "user@example.com",
		Name:    "Test User",
		Company: "Test Company",
	}, nil
}
