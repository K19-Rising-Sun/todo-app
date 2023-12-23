package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var schema = `
CREATE TABLE IF NOT EXISTS user (
    username VARCHAR(250) NOT NULL PRIMARY KEY, 
    password VARCHAR(250) NOT NULL
);
CREATE TABLE IF NOT EXISTS todo (
    id INTEGER PRIMARY KEY, 
    username TEXT NOT NULL, 
    category TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    state TEXT NOT NULL
)`

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

type Todo struct {
	Id          int    `db:"id"`
	Username    string `db:"username"`
	Category    string `db:"category"`
	Title       string `db:"title"`
	Description string `db:"description"`
	State       string `db:"state"`
}

func Init() (*sqlx.DB, error) {
	sqlite, err := sqlx.Connect("sqlite3", "todo.db")
	if err != nil {
		return nil, err
	}
	sqlite.MustExec(schema)
	return sqlite, nil
}
