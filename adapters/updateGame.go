package adapters

import (
	_ "github.com/jackc/pgx/stdlib"
)

func (s *dbGameService) UpdateGame(id int, players []string) error {
	var err error
	_, err = s.db.Exec("UPDATE game SET players = $1 WHERE id = $2", players, id)
	return err
}
