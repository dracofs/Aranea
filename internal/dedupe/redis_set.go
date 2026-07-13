package dedupe

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisSet struct {
	client* redis.Client
	contxt context.Context
}

func NewRedisSet() *RedisSet {
	rbd := redis.NewClient(&redis.Options{
		Addr: "localhost:6379"
	})

	return &RedisSet{
		client: rbd,
		contxt: context.Background()
	}

}

func (s *RedisSet) Seen(url string) (bool, error) {
	exists, err = s.client.SIsMember(s.contxt, "visited_urls", url).Result()
	
	if err != nil {
		return false, err
	}
	
	return exists, nil
}

func (s *RedisSet) Add(url string) error {
	return s.client.SAdd(s.contxt, "visited_urls", url).Err()
}
