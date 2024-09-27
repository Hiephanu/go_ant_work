package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService interface {
	// Health returns a map of health status information for Redis.
	// The keys and values in the map are Redis-specific.
	Health() map[string]string

	// Close terminates the Redis connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// Set a value in Redis with a given key and expiration time.
	Set(key string, value string, expiration time.Duration) error

	// Get a value from Redis by key.
	Get(key string) (string, error)

	// Delete a value from Redis by key.
	Delete(key string) error

	// Increment a value in Redis by key.
	Increment(key string) (int64, error)

	// Decrement a value in Redis by key.
	Decrement(key string) (int64, error)

	// Check if a key exists in Redis.
	Exists(key string) (bool, error)

	// Set a field in a hash stored at key.
	HashSet(key string, field string, data interface{}, expiration time.Duration) error

	// get a field in a hash stored at key
	HashGet(key string, field string) (interface{}, error)

	//check exits field in hash
	HashExists(key string, field string) (bool, error)

	HashIncrement(key string, field string) (int64, error)

	HashDecrement(key string, field string) (int64, error)
}

type redisService struct {
	client *redis.Client
	ctx    context.Context
}

var (
	redisAddr     = os.Getenv("REDIS_ADDR")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisDB       = 0 // default database index
	redisInstance *redisService
)

func NewRedisService() RedisService {
	// Reuse Connection
	if redisInstance != nil {
		return redisInstance
	}

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisAddr, redisPort),
		Password: redisPassword, // no password set
		DB:       redisDB,       // use default DB
	})

	redisInstance = &redisService{
		client: rdb,
		ctx:    context.Background(),
	}
	return redisInstance
}

// Health checks the health of the Redis connection.
func (r *redisService) Health() map[string]string {
	stats := make(map[string]string)

	// Ping the Redis server
	_, err := r.client.Ping(r.ctx).Result()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("redis down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"
	return stats
}

// Close terminates the Redis connection.
func (r *redisService) Close() error {
	log.Printf("Disconnected from Redis: %s:%s", redisAddr, redisPort)
	return r.client.Close()
}

// Set a value in Redis with a given key and expiration time.
func (r *redisService) Set(key string, value string, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Get a value from Redis by key.
func (r *redisService) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// Delete a value from Redis by key.
func (r *redisService) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Increment a value in Redis by key.
func (r *redisService) Increment(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// Decrement a value in Redis by key.
func (r *redisService) Decrement(key string) (int64, error) {
	return r.client.Decr(r.ctx, key).Result()
}

// Check if a key exists in Redis.
func (r *redisService) Exists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// HashSet sets a field in a hash stored at key.
func (r *redisService) HashSet(key string, field string, data interface{}, expiration time.Duration) error {
	// Convert data to a string if necessary
	value := fmt.Sprintf("%v", data)
	return r.client.HSet(r.ctx, key, field, value, expiration).Err()
}

func (r *redisService) HashGet(key string, field string) (interface{}, error) {
	value, err := r.client.HGet(r.ctx, key, field).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (r *redisService) HashExists(key string, feild string) (bool, error) {
	result, err := r.client.HExists(r.ctx, key, feild).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

func (r *redisService) HashIncrement(key string, field string) (int64, error) {
	result, err := r.client.HIncrBy(context.Background(), key, field, 1).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r *redisService) HashDecrement(key string, field string) (int64, error) {
	result, err := r.client.HIncrBy(context.Background(), key, field, -1).Result()

	if err != nil {
		return 0, nil
	}

	return result, nil
}
