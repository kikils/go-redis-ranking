package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kikils/go-redis-ranking/internal/domain"
	"github.com/kikils/go-redis-ranking/internal/infrastructure/redis"
	redis_repository "github.com/kikils/go-redis-ranking/internal/interface/redis"
	"github.com/kikils/go-redis-ranking/internal/usecase"
)

func main() {
	ctx := context.Background()
	redisCon, err := redis.NewRing([]string{fmt.Sprintf(
		"redis://%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)})
	if err != nil {
		log.Fatalf("redisCon: %v", err)
	}
	// mysqlCon, err := mysql.NewMySQLConnection()
	// if err != nil {
	// 	log.Fatalf("mysqlcon: %v", err)
	// }
	repository := redis_repository.NewRepository(redisCon)
	// repository := mysql_repository.NewRepository(mysqlCon)
	rankingInputPort := usecase.NewRankingUsecase(repository)
	result, err := rankingInputPort.GetRanking(ctx, 0, 100)
	if err != nil {
		log.Fatalf("GetRanking: %v", err)
	}
	resultLog(result)
}

func resultLog(result []*domain.UserRanking) {
	for _, v := range result {
		log.Println(v)
	}
}
