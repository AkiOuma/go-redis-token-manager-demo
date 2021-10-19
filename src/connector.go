package src

import "github.com/go-redis/redis/v8"

func ConnectRedis(config *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Url(),
		Password: "",
		DB:       0,
	})
}
