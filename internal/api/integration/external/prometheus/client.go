package prometheus

import "context"

type Client interface {
	Inc(ctx context.Context, name string, labels map[string]string) error
	Observe(ctx context.Context, name string, value float64, labels map[string]string) error
}

type client struct{}

func NewClient() Client {
	return &client{}
}

func (c *client) Inc(ctx context.Context, name string, labels map[string]string) error {
	return nil
}

func (c *client) Observe(ctx context.Context, name string, value float64, labels map[string]string) error {
	return nil
}
