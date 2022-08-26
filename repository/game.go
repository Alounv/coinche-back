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

var playerSchema = `
CREATE TABLE player (
	id serial PRIMARY KEY NOT NULL,
	name text NOT NULL,
	team text,
	gameid integer NOT NULL REFERENCES game(id),
	createdAt timestamp NOT NULL DEFAULT now(),
	initialOrder integer DEFAULT 0,
	cOrder integer DEFAULT 0,
	hand json NOT NULL DEFAULT '[]'
)`

var bidSchema = `
CREATE TABLE bid (
	id serial PRIMARY KEY NOT NULL,
	gameid integer NOT NULL REFERENCES game(id),
	value integer NOT NULL,
	player  text NOT NULL,
	color  text NOT NULL,
	coinche integer DEFAULT 0,
	pass integer DEFAULT 0
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

func (s *GameRepository) CreateBidsTableIfNeeded() error {
	_, err := s.db.Exec(bidSchema)
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

	err = gameRepository.CreateBidsTableIfNeeded()
	if err != nil {
		return &gameRepository, err
	}

	err = gameRepository.CreatePlayerTableIfNeeded()

	return &gameRepository, err
}
