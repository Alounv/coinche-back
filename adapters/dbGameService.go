package adapters

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

var gameSchema = `
CREATE TABLE game (
	id serial PRIMARY KEY NOT NULL,
	name text,
	createdAt timestamp NOT NULL DEFAULT now()
)`

type dbGameService struct {
	db *sqlx.DB
}

func (s *dbGameService) CreatePlayerTableIfNeeded() {
	s.db.Exec(gameSchema)
}

func NewDBGameService(dsn string) *dbGameService {
	db := sqlx.MustOpen("pgx", dsn)

	return NewDbGameServiceFromDb(db)
}

func NewDbGameServiceFromDb(db *sqlx.DB) *dbGameService {
	service := dbGameService{db}
	service.CreatePlayerTableIfNeeded()

	return &service
}
