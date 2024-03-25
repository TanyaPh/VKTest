package floodmonitoring

import (
	"context"
	"errors"
	"task/repository/postgres"
)

type FloodControl interface {
	Check(ctx context.Context, userID int64) (bool, error)
}

type FloodMonitoring struct {
	// usersChecks map[int64]int
	repo *postgres.Repository
	checksLimit int
}

func NewFloodMonitoring(repo *postgres.Repository, checksLimit int) *FloodMonitoring {
	return &FloodMonitoring {
		// usersChecks: make(map[int64]int),
		repo: repo,
		checksLimit: checksLimit,
	}
}

func (fm *FloodMonitoring) Check(ctx context.Context, userID int64) (bool, error) {
	select {
	case <-ctx.Done():
		return true, errors.New("The control time has ended")
	default:
		// checksCount, err := fm.repo.GetChecks(userID)
		checksCount, err := fm.repo.AddCheck(userID)
		return checksCount < fm.checksLimit, err
	}
}
