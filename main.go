package main

import (
	"coinche/app"
	"coinche/api"
	"coinche/utilities/env"
	"log"
	"os"

	_ "github.com/jackc/pgx/stdlib"
)

func main() {
	env.LoadEnv("")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")
	dbName := os.Getenv("DB_NAME")
	addr := os.Getenv("PORT")

	dsn := connectionInfo + " dbname=" + dbName
	gameService := app.NewGameService(dsn)

	router := api.SetupRouter(gameService)

	log.Print("Listening on ", addr)
	router.Run(addr)
}
