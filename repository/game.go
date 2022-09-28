package repository

import (
	"coinche/usecases"

	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

var gameSchema = `
CREATE TABLE game (
	id serial PRIMARY KEY NOT NULL,
	name text NOT NULL,
	createdAt timestamp NOT NULL DEFAULT now(),
	phase integer DEFAULT 0,
	deck json NOT NULL DEFAULT '[]'
)`

type GameRepository struct {
	usecases.GameRepositoryInterface
	db *sqlx.DB
}

func (s *GameRepository) CreateGameTableIfNeeded() error {
	_, err := s.db.Exec(gameSchema)
	return err
}

func (s *GameRepository) CreatePlayerTableIfNeeded() error {
	_, err := s.db.Exec(playerSchema)
	return err
}

func (s *GameRepository) CreateBidTableIfNeeded() error {
	_, err := s.db.Exec(bidSchema)
	return err
}

func (s *GameRepository) CreateTurnTableIfNeeded() error {
	_, err := s.db.Exec(turnSchema)
	return err
}

func (s *GameRepository) CreatePointTableIfNeeded() error {
	_, err := s.db.Exec(pointSchema)
	return err
}

func (s *GameRepository) CreateScoreTableIfNeeded() error {
	_, err := s.db.Exec(scoreSchema)
	return err
}
func NewGameRepository(dsn string) (*GameRepository, error) {
	db := sqlx.MustOpen("pgx", dsn)

	return NewGameRepositoryFromDb(db)
}

func NewGameRepositoryFromDb(db *sqlx.DB) (*GameRepository, error) {
	gameRepository := GameRepository{db: db}

	err := gameRepository.CreateGameTableIfNeeded()
	if err != nil {
		return &gameRepository, err
	}

	err = gameRepository.CreateBidTableIfNeeded()
	if err != nil {
		return &gameRepository, err
	}

	err = gameRepository.CreateTurnTableIfNeeded()
	if err != nil {
		return &gameRepository, err
	}

	err = gameRepository.CreatePointTableIfNeeded()
	if err != nil {
		return &gameRepository, err
	}

	err = gameRepository.CreateScoreTableIfNeeded()
	if err != nil {
		return &gameRepository, err
	}

	err = gameRepository.CreatePlayerTableIfNeeded()

	return &gameRepository, err
}
