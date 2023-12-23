package database

import (
	"github.com/jmoiron/sqlx"
)

func (todo *Todo) Create(db *sqlx.DB) error {
	_, err := db.NamedExec("INSERT INTO todo (username, category, title, description, state) VALUES (:username, :category, :title, :description, :state)", todo)
	return err
}

func (todo *Todo) Update(db *sqlx.DB) error {
	_, err := db.NamedExec(
		"UPDATE todo SET username = :username, category = :category, title = :title, description = :description, state = :state WHERE id = :id",
		todo,
	)
	return err
}

func GetTodos(db *sqlx.DB, username string) ([]Todo, error) {
	todos := []Todo{}
	err := db.Select(&todos, "SELECT * FROM todo WHERE username = ?", username)
	return todos, err
}

func DeleteTodo(db *sqlx.DB, id int) error {
	_, err := db.Exec("DELETE FROM todo WHERE id = ?", id)
	return err
}
