package floodmonitoring

import (
	"context"
	"task/repository/postgres"
)

type FloodControl interface {
	Check(ctx context.Context, userID int64) (bool, error)
}

type FloodMonitoring struct {
	repo *postgres.Repository
	checksLimit int
	timeLimit int
}

func NewFloodMonitoring(repo *postgres.Repository, checksLimit int, timeLimit int) *FloodMonitoring {
	return &FloodMonitoring {
		repo: repo,
		checksLimit: checksLimit,
		timeLimit: timeLimit,
	}
}

func (fm *FloodMonitoring) Check(ctx context.Context, userID int64) (bool, error) {
		checksCount, err := fm.repo.CountChecks(userID, fm.timeLimit)
		return checksCount < fm.checksLimit, err
}
