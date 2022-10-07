package main

import (
	"context"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

var REDIS *redis.Client

func parseRedisOptionsFromUri(uri string) *redis.Options {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	var db = strings.Trim(u.Path, "/")
	if db == "" {
		db = "0"
	}
	dbNum, err := strconv.Atoi(db)
	if err != nil {
		panic("Invalid Redis DB number")
	}

	password, _ := u.User.Password()

	return &redis.Options{
		Addr:     u.Host,
		Password: password,
		DB:       dbNum,
	}
}

func SetupRedis() {
	ctx := context.TODO()
	client := redis.NewClient(parseRedisOptionsFromUri(GetEnv("REDIS_CONNECTION")))
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}
	REDIS = client
}

func CloseRedis() {
	if REDIS != nil {
		log.Println("Closing redis connection")
		REDIS.Close()
	}
}
