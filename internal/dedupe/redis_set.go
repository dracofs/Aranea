package dedupe

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisSet struct {
	client *redis.Client
	contxt context.Context
}

func NewRedisSet(client *redis.Client) *RedisSet {
	return &RedisSet{
		client: client,
		contxt: context.Background(),
	}
}

func (s *RedisSet) Seen(url string) (bool, error) {
	exists, err := s.client.SIsMember(s.contxt, "visited_urls", url).Result()

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *RedisSet) Add(url string) (bool, error) {
	n, err := s.client.SAdd(s.contxt, "visited_urls", url).Result()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func (s *RedisSet) Clear() error {
	return s.client.Del(s.contxt, "visited_urls").Err()
}
