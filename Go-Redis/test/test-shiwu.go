package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var ctx = context.Background()
/*
在 Go-Redis 中，TxPipelined 和 TxPipeline 函数可以自动管理事务。
如果在执行管道中的命令时出现错误，这些函数会自动调用 DISCARD 命令来放弃事务。
 */
func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	defer  rdb.Close()
	key := "counter"

	// 使用 watch 来实现乐观锁
	err := rdb.Watch(ctx, func(tx *redis.Tx) error {
		n, err := tx.Get(ctx, key).Int()
		if err != nil {
			return err
		}

		// 事务操作
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// 递增操作
			pipe.Set(ctx, key, n+1, 0)
			return nil
		})
		return err
	}, key)

	if err != nil {
		fmt.Println("Transaction failed: ", err)
	} else {
		fmt.Println("Transaction success")
	}
}


func main2() {
	// 创建一个客户端对象
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",  // Redis 服务器地址
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	// 使用 Watch 方法来监视一个键，并在事务中操作这个键
	err := client.Watch(ctx, func(tx *redis.Tx) error {
		// 在事务开始之前，使用 Get 方法获取键 "count" 的当前值
		n, err := tx.Get(ctx, "count").Int()
		if err != nil && err != redis.Nil {
			return err
		}

		// 使用 Multi 方法开始一个事务
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// 在事务中，使用 Set 方法将键 "count" 的值加 1
			pipe.Set(ctx, "count", n+1, 0)
			return nil
		})
		return err
	}, "count")

	if err != nil {
		fmt.Println("transaction failed:", err)
		return
	} else {
		fmt.Println("Transaction success")
	}

	// 获取键 "count" 的新值
	val, err := client.Get(ctx, "count").Result()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("count:", val)

	// 输出: count: 1 (或者比 1 大的数字，取决于事务执行的次数)
}

