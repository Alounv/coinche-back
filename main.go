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
	gameRepositary := gamerepo.NewGameRepository(dsn)
	gameUsecases := usecases.NewGameUsecases(gameRepositary)

	router := api.SetupRouter(gameUsecases)

	log.Print("Listening on ", addr)
	err := router.Run(addr)
	if err != nil {
		panic(err)
	}
}
