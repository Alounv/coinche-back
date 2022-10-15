package repository

import "coinche/domain"

func (s *GameRepository) ListGames() ([]domain.Game, error) {
	var games []domain.Game
	gamesIDs := []int{}

	tx := s.db.MustBegin()

	err := tx.Select(&gamesIDs, "SELECT id FROM game WHERE id = root")
	if err != nil {
		return nil, err
	}

	for _, gameID := range gamesIDs {
		game, err := getGame(tx, gameID)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, tx.Commit()
}
