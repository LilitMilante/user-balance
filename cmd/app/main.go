package main

import (
	"log"
	"user-balance/bootstrap"
	"user-balance/domain/service"
	"user-balance/handlers"
	"user-balance/store"

	_ "github.com/lib/pq"
)

func main() {
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
	bs := service.NewBalanceService(s)
	srv := handlers.NewServer(bs)
	err = srv.Start(c.HttpPort)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
