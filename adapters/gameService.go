package adapters

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

var gameSchema = `
CREATE TABLE game (
	id serial PRIMARY KEY NOT NULL,
	name text
)`

type GameService struct {
	db *sqlx.DB
}

func (s *GameService) CreatePlayerTableIfNeeded() {
	s.db.Exec(gameSchema)
}

func NewGameService(dsn string) *GameService {
	db := sqlx.MustOpen("pgx", dsn)

	return NewGameServiceFromDb(db)
}

func NewGameServiceFromDb(db *sqlx.DB) *GameService {

	store := GameService{db}
	store.CreatePlayerTableIfNeeded()

	return &store
}
