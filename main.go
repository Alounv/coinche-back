package main

import (
	"coinche/api"
	gamerepo "coinche/repository/game"
	"coinche/usecases"
	"coinche/utilities/env"
	"log"
	"os"
)

func main() {
	env.LoadEnv("")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")
	dbName := os.Getenv("DB_NAME")
	addr := os.Getenv("PORT")

	dsn := connectionInfo + " dbname=" + dbName
	gameRepo := gamerepo.NewGameRepository(dsn)
	gameService := usecases.NewGameService(gameRepo)

	router := api.SetupRouter(gameService)

	log.Print("Listening on ", addr)
	err := router.Run(addr)
	if err != nil {
		panic(err)
	}
}
