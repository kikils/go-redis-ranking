package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/kikils/go-redis-ranking/internal/domain"
	"github.com/kikils/go-redis-ranking/internal/infrastructure/redis"
	redis_repository "github.com/kikils/go-redis-ranking/internal/interface/redis"
	"github.com/kikils/go-redis-ranking/internal/usecase"
	"golang.org/x/sync/errgroup"
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

	if err := repository.Del(ctx); err != nil {
		log.Fatalf("Del: %v", err)
	}
	var userIDs []string
	for i := 0; i < 1000; i++ {
		uuid := uuid.NewString()
		userIDs = append(userIDs, uuid)
	}

	group, _ := errgroup.WithContext(ctx)
	for _, userID := range userIDs {
		u := userID
		group.Go(func() (err error) {
			for i := 0; i < 100; i++ {
				wait, err := rand.Int(rand.Reader, big.NewInt(5))
				if err != nil {
					return err
				}
				time.Sleep(time.Duration(wait.Int64()) * time.Second)
				incr, err := rand.Int(rand.Reader, big.NewInt(10))
				if err != nil {
					return err
				}
				if err := rankingInputPort.AddScore(ctx, u, int(incr.Int64())); err != nil {
					return fmt.Errorf("AddScore: %v", err)
				}
			}
			return nil
		})
	}
	// group.Go(func() (err error) {
	// 	for i := 0; i < 10; i++ {
	// 		time.Sleep(100 * time.Millisecond)
	// 		result, err := rankingInputPort.GetRanking(ctx, 0, 100)
	// 		if err != nil {
	// 			return fmt.Errorf("GetRanking: %v", err)
	// 		}
	// 		log.Println(result)
	// 	}
	// 	return nil
	// })
	if err := group.Wait(); err != nil {
		log.Fatalf("error_group: %s", err)
	}
	var result []*domain.UserRanking
	Benchmark_execProcess(
		func() {
			result, err = rankingInputPort.GetRanking(ctx, 0, 100)
			if err != nil {
				log.Fatalf("GetRanking: %v", err)
			}
		},
	)
	resultLog(result)
}

func resultLog(result []*domain.UserRanking) {
	for _, v := range result {
		log.Println(v)
	}
}

func Benchmark_execProcess(f func()) {
	s := time.Now()
	f()
	fmt.Printf("process time: %s\n", time.Since(s))
}
