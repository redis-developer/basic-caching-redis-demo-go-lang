package main

import (
	"github.com/redis-developer/basic-caching-redis-demo-go-lang/api"
	"github.com/redis-developer/basic-caching-redis-demo-go-lang/config"
	"github.com/redis-developer/basic-caching-redis-demo-go-lang/controller"
	"github.com/redis-developer/basic-caching-redis-demo-go-lang/internal"
	"github.com/redis-developer/basic-caching-redis-demo-go-lang/redis"
	"log"
)

func main() {

	newConfig := config.NewConfig()

	newRedis := redis.NewRedis(newConfig.Redis)
	newServer := api.NewServer(newConfig.Api)

	controller.SetRedis(newRedis)

	internal.Waiting(newServer, newRedis)

	log.Println("Server exiting")

}
