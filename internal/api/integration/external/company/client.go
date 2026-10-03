package company

import (
	"context"

	usersservice "go-arch-template/internal/api/integration/external/usersservice"
)

type Client interface {
	ValidateCompany(ctx context.Context, companyID string) (bool, error)
	SyncCompany(ctx context.Context, companyID string) error
}

type client struct {
	users usersservice.Client
}

func NewClient(users usersservice.Client) Client {
	return &client{users: users}
}

func (c *client) ValidateCompany(ctx context.Context, companyID string) (bool, error) {
	if companyID == "" {
		return false, ErrInvalid
	}
	// Placeholder: validate via users-service in real wiring.
	_ = c.users
	return true, nil
}

func (c *client) SyncCompany(ctx context.Context, companyID string) error {
	if companyID == "" {
		return ErrInvalid
	}
	return nil
}
