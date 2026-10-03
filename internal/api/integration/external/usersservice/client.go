package usersservice

import "context"

type Client interface {
	GetUser(ctx context.Context, userID string) (*User, error)
	ValidateUser(ctx context.Context, userID string) (bool, error)
	SyncUser(ctx context.Context, userID string) error
}

type client struct {
	users map[string]*User
}

func NewClient() Client {
	return &client{
		users: make(map[string]*User),
	}
}

func (c *client) GetUser(ctx context.Context, userID string) (*User, error) {
	user, ok := c.users[userID]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (c *client) ValidateUser(ctx context.Context, userID string) (bool, error) {
	_, ok := c.users[userID]
	return ok, nil
}

func (c *client) SyncUser(ctx context.Context, userID string) error {
	return nil
}
