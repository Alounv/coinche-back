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
			INSERT INTO game (id, name, createdAt, phase, players)
			VALUES ($1, $2, $3, $4, string_to_array($5, ','))
			`,
			game.ID,
			game.Name,
			game.CreatedAt,
			game.Phase,
			strings.Join(game.Players, ","),
		)
	}
	err := tx.Commit()
	if err != nil {
		panic(err)
	}
}
