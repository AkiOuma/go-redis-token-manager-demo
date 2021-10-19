package src

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type TokenServer struct{}

func NewTokenServer() *TokenServer {
	return &TokenServer{}
}

// 模拟服务返回token，请求时长为0-10之间的一个随机数，单位是秒
//
// 用uuid来模拟一个具体的token的值
func (TokenServer) NewToken() string {
	rand.Seed(time.Now().Unix())
	waiting := rand.Intn(10)
	time.Sleep(time.Second * time.Duration(waiting))
	token := uuid.New().String()
	fmt.Printf("NewToken:%s  ||  Generated at %v, waiting for %v second\n", token, time.Now(), int(waiting))
	return token
}
