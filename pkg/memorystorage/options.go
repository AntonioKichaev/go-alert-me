package memorystorage

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

func SetStoreIntervalSecond(i int) Option {
	return func(storage *MemoryStorage) {
		storage.storeIntervalSecond = i
	}
}
