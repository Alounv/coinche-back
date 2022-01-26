package main

import (
	"coinche/adapters"
	"coinche/ports"
	"coinche/utilities/env"
	"log"
	"os"
)

func main() {
	env.LoadEnv("../.env")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")
	dbName := os.Getenv("DB_NAME")
	addr := os.Getenv("PORT")

	dsn := connectionInfo + " dbname=" + dbName
	gameService := adapters.NewGameService(dsn)

	router := ports.SetupRouter(gameService)

	log.Print("Listening on ", addr)
	router.Run(addr)
}
