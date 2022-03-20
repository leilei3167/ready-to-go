package main

import (
	"log"

	"github.com/go-redis/redis"
)

func main() {
	//ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set("k1", "m121212", 0).Err()
	if err != nil {
		log.Fatal(err)
	}

}
