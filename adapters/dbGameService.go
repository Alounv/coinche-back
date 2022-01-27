package adapters

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

type dbGameService struct {
	db *sqlx.DB
}

func (s *dbGameService) CreatePlayerTableIfNeeded() {
	s.db.Exec(gameSchema)
}

func newDBGameService(dsn string) *dbGameService {
	db := sqlx.MustOpen("pgx", dsn)

	return NewDbGameServiceFromDb(db)
}

func NewDbGameServiceFromDb(db *sqlx.DB) *dbGameService {
	service := dbGameService{db}
	service.CreatePlayerTableIfNeeded()

	return &service
}
