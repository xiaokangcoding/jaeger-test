package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// 创建一个发布者
	go Publisher(rdb)

	// 创建一个订阅者
	go Subscriber(rdb)

	// 这只是为了阻止主goroutine退出，实际的聊天应用程序可能会有其他的方式来处理这个问题
	select {}
}

func Publisher(rdb *redis.Client) {
	for {
		// 这里可以接收用户输入的消息，然后发布到聊天室
		var message string
		fmt.Print("Enter message: ")
		fmt.Scanln(&message)

		// 将消息发布到聊天室
		rdb.Publish(ctx, "chatroom", message)
	}
}

func Subscriber(rdb *redis.Client) {
	// 订阅聊天室的消息
	pubsub := rdb.Subscribe(ctx, "chatroom")
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Println(err.Error())
	}

	// 获取订阅的消息
	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Println("Received message: ", msg.Payload)
	}
}
