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
	phase integer DEFAULT 0,
	players text[]
)`

type GameRepository struct {
	usecases.GameRepositoryInterface
	db *sqlx.DB
}

func (s *GameRepository) CreatePlayerTableIfNeeded() error {
	_, err := s.db.Exec(gameSchema)
	return err
}

func NewGameRepository(dsn string) *GameRepository {
	db := sqlx.MustOpen("pgx", dsn)

	return NewGameRepositoryFromDb(db)
}

func NewGameRepositoryFromDb(db *sqlx.DB) *GameRepository {
	gameRepository := GameRepository{db: db}
	err := gameRepository.CreatePlayerTableIfNeeded()
	if err != nil {
		fmt.Println("No need to create the table.", err)
	}

	return &gameRepository
}
