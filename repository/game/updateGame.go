package gameRepo

import (
	_ "github.com/jackc/pgx/stdlib"
)

func (s *GameRepo) UpdateGame(id int, players []string) error {
	var err error
	_, err = s.db.Exec("UPDATE game SET players = $1 WHERE id = $2", players, id)
	return err
}
