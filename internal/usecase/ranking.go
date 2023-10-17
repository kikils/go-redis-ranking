package usecase

import (
	"context"

	"github.com/kikils/go-redis-ranking/internal/domain"
	"github.com/kikils/go-redis-ranking/internal/interface/repository"
)

type RankingInputPort interface {
	AddScore(ctx context.Context, userID string, incr int) error
	GetRanking(ctx context.Context, start, end int64) ([]*domain.UserRanking, error)
}

type RankingOutput struct {
}

type RankingInteractor struct {
	repository repository.Repository
}

func NewRankingUsecase(r repository.Repository) RankingInputPort {
	return &RankingInteractor{
		repository: r,
	}
}

func (i *RankingInteractor) AddScore(ctx context.Context, userID string, incr int) error {
	if err := i.repository.IncrScore(ctx, userID, incr); err != nil {
		return err
	}
	return nil
}

func (i *RankingInteractor) GetRanking(ctx context.Context, start, end int64) ([]*domain.UserRanking, error) {
	result, err := i.repository.GetRanking(ctx, start, end)
	if err != nil {
		return nil, err
	}
	return result, nil
}
