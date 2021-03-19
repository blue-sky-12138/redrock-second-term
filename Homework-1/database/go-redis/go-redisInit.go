package DBFirst

import (
	"github.com/go-redis/redis"
)

func LinkInit() error {
	client := redis.NewClient(&redis.Options{})
	_, err := client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
