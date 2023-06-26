package memorystorage

import (
	"encoding/json"
	"github.com/antoniokichaev/go-alert-me/internal/logger"
	"go.uber.org/zap"
	"io"
	"os"
	"sync"
	"time"
)

type saveFormat struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MemoryStorage struct {
	store               map[string]string
	storeIntervalSecond int
	pathToSaveLoad      string
	mu                  sync.RWMutex
}

func NewMemoryStorage(opts ...Option) (*MemoryStorage, error) {
	m := &MemoryStorage{
		store: make(map[string]string),
	}
	for _, opt := range opts {
		opt(m)
	}
	if m.storeIntervalSecond > 0 {
		go func() {
			t := time.NewTicker(time.Second * time.Duration(m.storeIntervalSecond))
			for range t.C {
				m.saveOnDisk()
			}
		}()
	}
	return m, nil
}

func (m *MemoryStorage) Set(name, value string) error {
	m.mu.Lock()
	m.store[name] = value
	m.mu.Unlock()
	m.DoSave()
	return nil
}

func (m *MemoryStorage) Get(name string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if val, ok := m.store[name]; ok {
		return val, nil
	}
	return "", ErrorNotExistMetric
}

func (m *MemoryStorage) GetDump() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]string, len(m.store))
	for key, val := range m.store {
		result[key] = val
	}
	return result
}
func (m *MemoryStorage) LoadFromDisk() error {
	file, err := os.Open(m.pathToSaveLoad)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(file)
	m.mu.Lock()
	defer m.mu.Unlock()
	for {
		el := saveFormat{}
		err = dec.Decode(&el)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		m.store[el.Key] = el.Value

	}
	return nil
}

func (m *MemoryStorage) saveOnDisk() {
	f, _ := os.OpenFile(m.pathToSaveLoad, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	enc := json.NewEncoder(f)
	m.mu.RLock()
	defer m.mu.RUnlock()
	for key, val := range m.store {
		el := &saveFormat{Key: key, Value: val}
		err := enc.Encode(el)
		if err != nil {
			logger.Log.Error("encoding", zap.Error(err))
			continue
		}
	}

}

func (m *MemoryStorage) DoSave() {
	if m.storeIntervalSecond == 0 {
		m.saveOnDisk()
	}
}
