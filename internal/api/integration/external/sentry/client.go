package sentry

import "context"

type Client interface {
	Capture(ctx context.Context, event Event) error
}

type client struct{}

func NewClient() Client {
	return &client{}
}

func (c *client) Capture(ctx context.Context, event Event) error {
	return nil
}
