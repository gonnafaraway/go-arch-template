package service

type API struct{}

func NewAPI() *API {
	return &API{}
}

func (a *API) Run() error {
	return nil
}

type Domain

func PrepareAPIService(envs Env) (*API, error) {}
