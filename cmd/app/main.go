package main

import (
	_ "github.com/lib/pq"
	"log"
	"user-balance/api"
	"user-balance/bootstrap"
	"user-balance/domain/service"
	"user-balance/infrastructure/redisdb"
	"user-balance/infrastructure/store"
)

func main() {
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	c, err := bootstrap.NewConfig()
	if err != nil {
		log.Fatal("config:", err)
	}

	db, err := bootstrap.ConnectDB(c)
	if err != nil {
		log.Fatal("connect Db:", err)
	}

	defer db.Close()

	rdb := redisdb.NewRedisStore()
	defer rdb.Close()

	s := store.NewStore(db, rdb)

	bs := service.NewBalanceService(s, c.ApiKey)

	srv := api.NewServer(bs)
	err = srv.Start(c.HttpPort)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
