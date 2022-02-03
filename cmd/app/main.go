package main

import (
	"context"
	"log"
	"user-balance/bootstrap"
	"user-balance/domain/service"
	"user-balance/handlers"
	"user-balance/infrastructure/redisdb"
	"user-balance/store"

	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := bootstrap.NewConfig()
	if err != nil {
		log.Fatal("config:", err)
	}

	db, err := bootstrap.ConnectDB(c)
	if err != nil {
		log.Fatal("connect Db:", err)
	}

	defer db.Close()

	s := store.NewStore(db)
	rdb := redisdb.NewRedisStore(ctx)
	//
	//err = rdb.Set(ctx, "Hello", "World", 0).Err()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//val, err := rdb.Get(ctx, "Hello").Result()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println("Hello", val)
	//

	bs := service.NewBalanceService(s, rdb, c.ApiKey)
	srv := handlers.NewServer(bs)
	err = srv.Start(c.HttpPort)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
