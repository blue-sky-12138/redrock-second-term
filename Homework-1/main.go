package main

import (
	DBSecond "Homework-1/database/redigo"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func main() { //这里为redigo的封装使用
	if err := DBSecond.LinkInit(); err != nil {
		return
	}
	poolExt := []DBSecond.PoolExt{ //自定义连接设置
		DBSecond.WithMaxIdle(40),
		DBSecond.WithMaxActive(10),
	}
	client := DBSecond.NewClient(poolExt...)

	client.Del("key1")

	client.ZAdd("key1", 1, "m1", 2, "m2", 3, "m3", 4, "m4")

	temCliceV, _ := client.ZRangeAll("key1")
	for _, temv := range temCliceV {
		v, _ := redis.String(temv, nil)
		fmt.Printf("%s ", v)
	}
	fmt.Printf("\n")

	client.ZRem("key1", "m2", "m3")

	temCliceV, _ = client.ZRangeAll("key1")
	for _, temv := range temCliceV {
		v, _ := redis.String(temv, nil)
		fmt.Printf("%s ", v)
	}
	fmt.Printf("\n")

}

//func main() { //这里为go-redis的封装
//	if err := DBFirst.LinkInit(); err != nil {
//		fmt.Println(err)
//		return
//	}
//	options := []DBFirst.DialOption{ //自定义连接设置
//		DBFirst.WithPoolSize(20),
//	}
//	client := DBFirst.NewClient(options)
//	client.Set("key1", "v1", 5*time.Minute)
//	v, err := client.Get("key1").Result()
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(v)
//}
