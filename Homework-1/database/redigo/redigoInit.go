package DBSecond

import (
	"github.com/gomodule/redigo/redis"
	"log"
)

func LinkInit() error {
	rdb, err := redis.Dial("tcp", "localhost:6379")
	defer rdb.Close()
	if err != nil {
		log.Println("LinkFail:", err)
		return err
	}
	return nil
}
