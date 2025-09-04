package main

const {
	dbDriver = "postgres"
	dbSource = "postgresql:.."
	serverAddress = "0.0.0.0:8080"
}

func main(){
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil{
		log.Fatal("Cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err := server.Start(serverAddress)
	if err != nil{
		log.Fatal("Cannot start server: ", err)
	}
}