package storage

type Model interface {
}

type Storage struct {
}

func PrepareStorage() (*Storage, error) {
	return &Storage{}, nil
}
