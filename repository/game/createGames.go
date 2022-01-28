package gameRepo

import (
	"coinche/domain"
	"strings"

	_ "github.com/jackc/pgx/stdlib"
)

func (s *GameRepo) CreateGames(games []domain.Game) {
	tx := s.db.MustBegin()
	for _, game := range games {
		tx.MustExec(
			`INSERT INTO game (id, name, createdAt, players)
			VALUES ($1, $2, $3, string_to_array($4, ','))`,
			game.Id, game.Name, game.CreatedAt, strings.Join(game.Players, ","),
		)
	}
	tx.Commit()
}
