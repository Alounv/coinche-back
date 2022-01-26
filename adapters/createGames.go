package adapters

import (
	"coinche/app"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

func (s *GameService) CreateGames(games []app.Game) {
	tx := s.db.MustBegin()
	for _, game := range games {
		_, err := tx.Exec("INSERT INTO game (name) VALUES ($1)", game.Name)
		if err != nil {
			fmt.Println(err)
		}
	}
	tx.Commit()
}
