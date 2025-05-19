package main

import (
	"github.com/lekchan000/isekai-shop-api/config"
	"github.com/lekchan000/isekai-shop-api/database"
	"github.com/lekchan000/isekai-shop-api/server"
)

func main() {
	conf := config.ConfigGetting()
	db := database.NewPostgresDatabase(conf.Database)
	server := server.NewEchoServer(conf, db)

	server.Start()
}
