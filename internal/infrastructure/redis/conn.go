package redis

import (
	"fmt"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func NewRing(urls []string) (*redis.Ring, error) {
	var (
		options = make(map[string]*redis.Options)
		addrs   = make(map[string]string)
	)
	for i, url := range urls {
		opt, err := redis.ParseURL(url)
		if err != nil {
			return nil, err
		}
		options[opt.Addr] = opt
		addrs[fmt.Sprintf("shard_%d", i)] = opt.Addr
	}
	rdb := redis.NewRing(&redis.RingOptions{
		NewClient: func(opt *redis.Options) *redis.Client {
			return redis.NewClient(options[opt.Addr])
		},
	})
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, err
	}
	rdb.SetAddrs(addrs)
	return rdb, nil
}
