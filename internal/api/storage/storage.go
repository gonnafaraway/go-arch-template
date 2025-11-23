package storage

type Storage struct {
	// Можно добавить реальные клиенты для БД, Kafka и т.д.
}

func PrepareStorage(env interface{}) (*Storage, error) {
	// Здесь можно инициализировать реальные подключения
	// Пока используем моки
	return &Storage{}, nil
}
