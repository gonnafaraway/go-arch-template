package oauth

import "context"

type Client interface {
	ValidateToken(ctx context.Context, token string) (*TokenInfo, error)
	GetUserInfo(ctx context.Context, token string) (*UserInfo, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
}

type client struct {
	tokens map[string]*TokenInfo
}

func NewClient() Client {
	return &client{
		tokens: make(map[string]*TokenInfo),
	}
}

func (c *client) ValidateToken(ctx context.Context, token string) (*TokenInfo, error) {
	if token == "" {
		return &TokenInfo{Valid: false}, nil
	}
	info, ok := c.tokens[token]
	if !ok {
		return &TokenInfo{Valid: true, UserID: "user_1"}, nil
	}
	return info, nil
}

func (c *client) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	info, err := c.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !info.Valid {
		return nil, ErrInvalidToken
	}
	return &UserInfo{
		UserID:  info.UserID,
		Email:   "user@example.com",
		Name:    "Test User",
		Company: "Test Company",
	}, nil
}

func (c *client) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	return &TokenResponse{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
		ExpiresIn:    3600,
	}, nil
}
