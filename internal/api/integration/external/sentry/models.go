package sentry

type Event struct {
	Message string            `json:"message"`
	Level   string            `json:"level"`
	Tags    map[string]string `json:"tags"`
}
