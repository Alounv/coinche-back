package gameRepo

import (
	"github.com/jmoiron/sqlx"
)

var gameSchema = `
CREATE TABLE game (
	id serial PRIMARY KEY NOT NULL,
	name text,
	createdAt timestamp NOT NULL DEFAULT now(),
	players text[]
)`

type GameRepo struct {
	db *sqlx.DB
}

func (s *GameRepo) CreatePlayerTableIfNeeded() {
	s.db.Exec(gameSchema)
}

func NewGameRepo(dsn string) *GameRepo {
	db := sqlx.MustOpen("pgx", dsn)

	return NewGameRepoFromDb(db)
}

func NewGameRepoFromDb(db *sqlx.DB) *GameRepo {
	service := GameRepo{db}
	service.CreatePlayerTableIfNeeded()

	return &service
}
