package service

import (
	"todo-app/database"

	"github.com/jmoiron/sqlx"
)

func CreateTodo(db *sqlx.DB, todo *database.Todo) error {
	_, err := db.NamedExec("INSERT INTO todo (username, category, title, description, state) VALUES (:username, :category, :title, :description, :state)", todo)
	return err
}

func GetTodos(db *sqlx.DB, username string) ([]database.Todo, error) {
	todos := []database.Todo{}
	err := db.Select(&todos, "SELECT * FROM todo WHERE username = ?", username)
	return todos, err
}

func UpdateTodo(db *sqlx.DB, new_todo *database.Todo) error {
	_, err := db.NamedExec(
		"UPDATE todo SET username = :username, category = :category, title = :title, description = :description, state = :state WHERE id = :id",
		new_todo,
	)
	return err
}

func DeleteTodo(db *sqlx.DB, id int) error {
	_, err := db.Exec("DELETE FROM todo WHERE id = ?", id)
	return err
}
