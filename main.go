package main

import (
	"database/sql"
	"log"

	"github.com/Adebobola01/Simple-bank---GO/api"
	db "github.com/Adebobola01/Simple-bank---GO/db/sqlc"
	"github.com/Adebobola01/Simple-bank---GO/util"
	_ "github.com/lib/pq"
)



func main(){
	config, err := util.LoadConfig(".")
	if err != nil{
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil{
		log.Fatal("Cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil{
		log.Fatal("Cannot start server: ", err)
	}
}