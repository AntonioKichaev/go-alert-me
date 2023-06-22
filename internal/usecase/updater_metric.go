package usecase

import (
	"github.com/antoniokichaev/go-alert-me/internal/entity"
)

type UpdaterUseCase struct {
	repo UpdaterRepo
}

func NewUpdater(repo UpdaterRepo) *UpdaterUseCase {
	return &UpdaterUseCase{repo: repo}
}

func (u *UpdaterUseCase) AddCounter(name string, value any) error {
	c, err := entity.NewCounter(name, value)
	if err != nil {
		return err
	}
	return u.repo.AddCounter(c)
}

func (u *UpdaterUseCase) SetGauge(name string, value any) error {
	g, err := entity.NewGauge(name, value)
	if err != nil {
		return err
	}
	return u.repo.SetGauge(g)
}
