package usecase

import (
	metrics2 "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
)

type UpdaterUseCase struct {
	repo UpdaterRepo
}

func NewUpdater(repo UpdaterRepo) *UpdaterUseCase {
	return &UpdaterUseCase{repo: repo}
}

func (u *UpdaterUseCase) AddCounter(name string, value any) (*metrics2.Counter, error) {
	c, err := metrics2.NewCounter(name, value)
	if err != nil {
		return nil, err
	}
	c, err = u.repo.AddCounter(c)
	return c, err
}

func (u *UpdaterUseCase) SetGauge(name string, value any) (*metrics2.Gauge, error) {
	g, err := metrics2.NewGauge(name, value)
	if err != nil {
		return nil, err
	}
	g, err = u.repo.SetGauge(g)
	return g, err
}
