package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client *redis.Client
	contxt context.Context
}

func NewRedisQueue(client *redis.Client) *RedisQueue {
	return &RedisQueue{
		client: client,
		contxt: context.Background(),
	}
}

func (q *RedisQueue) Push(url string) error {
	return q.client.LPush(q.contxt, "crawler_queue", url).Err()
}

func (q *RedisQueue) Pop() (string, error) {
	result, err := q.client.BRPop(q.contxt, 0, "crawler_queue").Result()

	if err != nil {
		return "", err
	}

	return result[1], nil
}

func (q *RedisQueue) Len() (int64, error) {
	return q.client.LLen(q.contxt, "crawler_queue").Result()
}

func (q *RedisQueue) Clear() error {
	return q.client.Del(q.contxt, "crawler_queue").Err()
}
