package queue

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sahilrana27582/go-task-queue/models"
)

const RedisTaskQueueKey = "task_queue"

type RedisQueue struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewRedisQueue(addr string) *RedisQueue {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		panic("failed to connect to Redis: " + err.Error())
	}

	return &RedisQueue{
		redisClient: client,
		ctx:         context.Background(),
	}
}

func (rq *RedisQueue) Enqueue(task *models.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return rq.redisClient.RPush(rq.ctx, RedisTaskQueueKey, data).Err()
}

func (rq *RedisQueue) Dequeue(timeout time.Duration) (*models.Task, error) {
	result, err := rq.redisClient.BLPop(rq.ctx, timeout, RedisTaskQueueKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("redis nil")
		}
		return nil, err
	}

	if len(result) < 2 {
		return nil, errors.New("invalid task data")
	}

	var task models.Task
	if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
		return nil, err
	}

	if err := rq.redisClient.LRem(rq.ctx, RedisTaskQueueKey, 0, result[1]).Err(); err != nil {
		return nil, err
	}

	task.Status = models.StatusInProgress

	return &task, nil
}

func (rq *RedisQueue) Length() (int64, error) {
	return rq.redisClient.LLen(rq.ctx, RedisTaskQueueKey).Result()
}

func (rq *RedisQueue) Ping() (string, error) {
	status, err := rq.redisClient.Ping(rq.ctx).Result()
	if err != nil {
		return "", err
	}
	return status, nil
}

func (rq *RedisQueue) Close() error {
	if rq.redisClient != nil {
		return rq.redisClient.Close()
	}
	return nil
}
