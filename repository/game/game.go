package gamerepo

import (
	"coinche/usecases"
	"fmt"

	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

var gameSchema = `
CREATE TABLE game (
	id serial PRIMARY KEY NOT NULL,
	name text,
	createdAt timestamp NOT NULL DEFAULT now(),
	players text[]
)`

type GameRepositary struct {
	usecases.GameRepositoryInterface
	db *sqlx.DB
}

func (s *GameRepositary) CreatePlayerTableIfNeeded() {
	_, err := s.db.Exec(gameSchema)
	if err != nil {
		fmt.Print(err)
	}
}

func NewGameRepository(dsn string) *GameRepositary {
	db := sqlx.MustOpen("pgx", dsn)

	return NewGameRepositaryFromDb(db)
}

func NewGameRepositaryFromDb(db *sqlx.DB) *GameRepositary {
	gameRepositary := GameRepositary{db: db}
	gameRepositary.CreatePlayerTableIfNeeded()

	return &gameRepositary
}
