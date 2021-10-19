package main

import (
	"fmt"
	"optimistic-lock/src"
	"sync"
	"time"
)

func main() {
	db := src.ConnectRedis(src.NewConfig(
		"127.0.0.1",
		"6379",
	))
	looper := 100
	server := src.NewTokenServer()
	wg := sync.WaitGroup{}
	wg.Add(looper)
	for i := 0; i < looper; i++ {
		if i%4 == 0 {
			time.Sleep(time.Second)
		}
		go func() {
			tokenManager := src.NewTokenManager(db, server)
			token, _ := tokenManager.GetToken()
			fmt.Println(token)
			defer wg.Done()
		}()
	}
	wg.Wait()
}
