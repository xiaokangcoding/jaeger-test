package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:	  "localhost:6379",
		Password: "", // no password set
		DB:		  0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		log.Fatal(err.Error())
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}