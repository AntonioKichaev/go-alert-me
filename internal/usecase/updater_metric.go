package usecase

import (
	"github.com/antoniokichaev/go-alert-me/pkg/metrics"
)

type UpdaterUseCase struct {
	repo UpdaterRepo
}

func NewUpdater(repo UpdaterRepo) *UpdaterUseCase {
	return &UpdaterUseCase{repo: repo}
}

func (u *UpdaterUseCase) AddCounter(name string, value any) (*metrics.Counter, error) {
	c, err := metrics.NewCounter(name, value)
	if err != nil {
		return nil, err
	}
	c, err = u.repo.AddCounter(c)
	return c, err
}

func (u *UpdaterUseCase) SetGauge(name string, value any) (*metrics.Gauge, error) {
	g, err := metrics.NewGauge(name, value)
	if err != nil {
		return nil, err
	}
	g, err = u.repo.SetGauge(g)
	return g, err
}
