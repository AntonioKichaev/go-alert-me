package memorystorage

import "go.uber.org/zap"

type Option func(storage *MemoryStorage)

func SetPathToSaveLoad(path string) Option {
	return func(storage *MemoryStorage) {
		storage.pathToSaveLoad = path
	}
}
func WithRestore(b bool) Option {
	return func(storage *MemoryStorage) {
		if b && storage.pathToSaveLoad != "" {
			err := storage.LoadFromDisk()
			if err != nil {
				return
			}
		}

	}
}
func WithLogger(logger *zap.Logger) Option {
	return func(storage *MemoryStorage) {
		storage.logger = logger
	}
}

func SetStoreIntervalSecond(i int) Option {
	return func(storage *MemoryStorage) {
		storage.storeIntervalSecond = i
	}
}
