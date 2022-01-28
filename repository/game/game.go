package gamerepo

import (
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

type GameRepo struct {
	db *sqlx.DB
}

func (s *GameRepo) CreatePlayerTableIfNeeded() {
	_, err := s.db.Exec(gameSchema)
	if err != nil {
		fmt.Print(err)
	}
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
