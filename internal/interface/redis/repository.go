package redis

import (
	"context"

	"github.com/kikils/go-redis-ranking/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	conn *redis.Ring
}

func NewRepository(conn *redis.Ring) *Repository {
	return &Repository{
		conn: conn,
	}
}

const (
	Key = "ranking"
)

func (r *Repository) IncrScore(ctx context.Context, userID string, incr int) error {
	conn := r.conn

	if err := conn.ZIncrBy(ctx, Key, float64(incr), userID).Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetRanking(ctx context.Context, start, end int64) ([]*domain.UserRanking, error) {
	conn := r.conn

	result, err := conn.ZRevRangeWithScores(ctx, Key, start, end).Result()
	if err != nil {
		return nil, err
	}
	userRankings := make([]*domain.UserRanking, len(result))
	for i, v := range result {
		userRankings[i] = &domain.UserRanking{
			UserID: v.Member.(string),
			Rank:   int(start) + i,
			Score:  int(v.Score),
		}
	}

	return userRankings, nil
}

func (r *Repository) Del(ctx context.Context) error {
	conn := r.conn

	if err := conn.Del(ctx, Key).Err(); err != nil {
		return err
	}

	return nil
}
