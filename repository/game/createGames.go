package gamerepo

import (
	"coinche/domain"
	"strings"
)

func (s *GameRepository) CreateGames(games []domain.Game) {
	tx := s.db.MustBegin()
	for _, game := range games {
		tx.MustExec(
			`
			INSERT INTO game (id, name, createdAt, players)
			VALUES ($1, $2, $3, string_to_array($4, ','))
			`,
			game.ID,
			game.Name,
			game.CreatedAt,
			strings.Join(game.Players, ","),
		)
	}
	err := tx.Commit()
	if err != nil {
		panic(err)
	}
}
