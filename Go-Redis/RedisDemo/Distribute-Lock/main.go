package main

import (
	"errors"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	logger "github.com/sirupsen/logrus"
	"time"
)

const mutexname = "my-global-mutex"

var rs * redsync.Redsync

func main() {
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs = redsync.New(pool)

	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.
	mutex := rs.NewMutex(mutexname)

	// Obtain a lock for our given mutex. After this is successful, no one else
	// can obtain the same lock (the same mutex name) until we unlock it.
	if err := mutex.Lock(); err != nil {
		logger.Error(err)
		return
	}
	// Release the lock so other processes or threads can obtain a lock.
	defer func() {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			logger.Errorf("Failed to release lock: %v", err)
		} else {
			logger.Info("Lock released")
		}
	}()

	logger.Info("Lock acquired, start doing work...")
	// 假设我们正在处理一个需要长时间的任务
	// 在另一goroutine中尝试延期锁
	errChan := make(chan error)
	go func() {
		errChan <- Retry(mutex)
	}()

//	go Work()
	// 假设我们正在处理一个需要长时间的任务
	select {
	case <-time.After(20 * time.Second):
		logger.Info("Work completed")
	case err := <-errChan:
		if err != nil {
			logger.Errorf("Failed to extend lock: %v", err)
			// 如果延期锁失败，停止工作
			return
		}
	}
}

func Retry(mutex *redsync.Mutex) error {
	for i := 0; i < 5; i++ {
		// Do some part of the work...
		time.Sleep(3 * time.Second)

		// Try to extend the lock.
		if ok, err := mutex.Extend(); !ok || err != nil {
			logger.Errorf("Failed to extend lock: %v", err)
			// If we failed to extend the lock, we should stop doing the work.
		//	return
		} else {
			logger.Info("Lock extended")
			return nil
		}
	}
	return errors.New("Failed to extend lock")
}

//func Work()  {
//	time.Sleep(time.Second *2)
//	mutex := rs.NewMutex(mutexname)
//
//	// Obtain a lock for our given mutex. After this is successful, no one else
//	// can obtain the same lock (the same mutex name) until we unlock it.
//	if err := mutex.Lock(); err != nil {
//		logger.Error(err)
//		return
//	}
//	time.Sleep(time.Second *10)
//}