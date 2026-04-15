package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client* redis.Client
	contxt context.Context
}

func NewRedisQueue() *RedisQueue {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &RedisQueue{
		client: rdb,
		contxt: context.Background(),
	}
}

func (q *RedisQueue) Push (url string) error {
	return q.client.LPush(q.contxt, "crawler_queue", url).Err()
}

func (q *RedisQueue) Pop() (string, error) {
	result, err := q.client.BRPop(q.contxt, 0, "crawler_queue").Result()

	if err != nil {
		return "", err
	}

	return result[1], nil
}