package env

import "os"

type Env struct {
	HTTPPort string
}

func PrepareEnv() (*Env, error) {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	return &Env{
		HTTPPort: port,
	}, nil
}
