package database

import (
	"net/http"
)

//go:generate mockery --name StorageStatus
type StorageStatus interface {
	Ping() error
}

type HadlerDB struct {
	db StorageStatus
}

func (h *HadlerDB) GetPing(w http.ResponseWriter, r *http.Request) {
	err := h.db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func NewHandlers(db StorageStatus) *HadlerDB {
	return &HadlerDB{db: db}
}
