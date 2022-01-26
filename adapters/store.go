package adapters

import (
	"coinche/app"
	"fmt"

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

func (s *GameService) CreateGames(games []app.Game) {
	tx := s.db.MustBegin()
	for _, game := range games {
		_, err := tx.Exec("INSERT INTO game (name) VALUES ($1)", game.Name)
		if err != nil {
			fmt.Println(err)
		}
	}
	tx.Commit()
}

func (s *GameService) CreateAGame(name string) int {
	var id int
	err := s.db.QueryRow("INSERT INTO game (name) VALUES ($1) RETURNING id", name).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}

func (s *GameService) GetAGame(id int) app.Game {
	var game app.Game
	err := s.db.Get(&game, "SELECT * FROM game WHERE id=$1", id)

	if err != nil {
		fmt.Println(err)
	}

	return game
}

func (s *GameService) GetAllGames() []app.Game {
	var games []app.Game
	err := s.db.Select(&games, "SELECT * FROM game ")

	if err != nil {
		fmt.Println(err)
	}

	return games
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
