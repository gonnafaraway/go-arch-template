package service

type API struct{}

func NewAPI() *API {
	return &API{}
}

func (a *API) Run() error {
	return nil
}

func PrepareAPIService(envs Env) (*API, error) {}
