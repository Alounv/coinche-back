package gamerepo

import (
	"coinche/domain"
	"strings"
)

func (s *GameRepository) CreateGame(game domain.Game) int {
	var id int

	err := s.db.QueryRow(
		`
		INSERT INTO game (name, players) 
		VALUES ($1, string_to_array($2, ',')) 
		RETURNING id
		`,
		game.Name,
		strings.Join(game.Players, ","),
	).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}
