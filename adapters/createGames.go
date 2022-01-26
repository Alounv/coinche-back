package adapters

import (
	"coinche/app"

	_ "github.com/jackc/pgx/stdlib"
)

func (s *GameService) CreateGames(games []app.Game) {
	tx := s.db.MustBegin()
	for _, game := range games {
		tx.MustExec(
			`INSERT INTO game (id, name, createdAt)
			VALUES ($1, $2, $3)`,
			game.Id, game.Name, game.CreatedAt,
		)
	}
	tx.Commit()
}
