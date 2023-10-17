package repository

import (
	"context"

	"github.com/kikils/go-redis-ranking/internal/domain"
)

type Repository interface {
	IncrScore(ctx context.Context, userID string, incr int) error
	GetRanking(ctx context.Context, start, end int64) ([]*domain.UserRanking, error)
}
