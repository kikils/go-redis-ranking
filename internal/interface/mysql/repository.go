package mysql

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/kikils/go-redis-ranking/internal/domain"
	"gorm.io/gorm"
)

type Repository struct {
	conn *gorm.DB
}

func NewRepository(conn *gorm.DB) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (r *Repository) IncrScore(ctx context.Context, userID string, incr int) error {
	u := &domain.UserScore{
		ID:    userID,
		Score: incr,
	}
	result := r.conn.WithContext(ctx).Table("user_scores").Where("id=?", userID).First(u)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("IncrScore: %s", result.Error)
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if err := r.conn.WithContext(ctx).Table("user_scores").Create(u).Error; err != nil {
			return fmt.Errorf("IncrScore: %s", err)
		}
	}
	u.Score = incr
	if err := r.conn.WithContext(ctx).Table("user_scores").Where("id=?", userID).Update("score", u.Score).Error; err != nil {
		log.Printf("IncrScore: %s", err)
		return err
	}
	return nil
}

func (r *Repository) GetRanking(ctx context.Context, start, end int64) ([]*domain.UserRanking, error) {
	var scores []*domain.UserScore
	if err := r.conn.WithContext(ctx).Table("user_scores").Limit(int(end)).Offset(int(start)).Find(&scores).Error; err != nil {
		return nil, err
	}
	var rankings []*domain.UserRanking
	for i, v := range scores {
		rankings = append(rankings, &domain.UserRanking{
			UserID: v.ID,
			Rank:   int(start) + i,
			Score:  v.Score,
		})
	}
	return rankings, nil
}

func (r *Repository) Del(ctx context.Context) error {
	if err := r.conn.Exec("TRUNCATE user_scores").Error; err != nil {
		return err
	}
	return nil
}
