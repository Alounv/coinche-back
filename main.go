package main

import (
	"coinche/api"
	repository "coinche/repository"
	"coinche/usecases"
	"coinche/utilities"
	"fmt"
	"os"
)

func main() {
	utilities.LoadEnv("")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")
	dbName := os.Getenv("DB_NAME")
	addr := os.Getenv("PORT")
	authorizedOrigin := os.Getenv("AUTHORIZED_ORIGIN")

	dsn := connectionInfo + " dbname=" + dbName
	gameRepository, err := repository.NewGameRepository(dsn)
	utilities.PanicIfErr(err)
	gameUsecases := usecases.NewGameUsecases(gameRepository)

	router, _ := api.SetupRouter(gameUsecases, []string{authorizedOrigin})

	fmt.Println("Listening on ", addr)
	err = router.Run(addr)
	utilities.PanicIfErr(err)
}
