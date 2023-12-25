package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (todo *Todo) Create(db *sqlx.DB) (Todo, error) {
	result, err := db.NamedExec("INSERT INTO todo (username, category, title, description, is_done) VALUES (:username, :category, :title, :description, :is_done)", todo)
	if err != nil {
		return Todo{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Todo{}, err
	}

	created_todo := Todo{}
	if err := db.Get(&created_todo, "SELECT * FROM todo WHERE id = ?", id); err != nil {
		return Todo{}, err
	}

	return created_todo, nil
}

func (todo *Todo) Update(db *sqlx.DB) (Todo, error) {
	result, err := db.NamedExec(
		"UPDATE todo SET username = :username, category = :category, title = :title, description = :description, is_done = :is_done WHERE id = :id",
		todo,
	)
	if err != nil {
		return Todo{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Todo{}, err
	}

	updated_todo := Todo{}
	if err := db.Get(&updated_todo, "SELECT * FROM todo WHERE id = ?", id); err != nil {
		return Todo{}, err
	}

	return updated_todo, nil
}

func QueryTodo(db *sqlx.DB, username string, category string, title string) ([]Todo, error) {
	todos := []Todo{}
	if category != "" && title != "" {
		err := db.Select(&todos, `SELECT * FROM todo WHERE username = ? AND id IN(
            SELECT rowid FROM todo_search WHERE todo_search MATCH ('category:' || ? || ' AND ' || ' title:' || ?)
        )`, username, category, title)
		return todos, err
	}
	if category != "" {
		err := db.Select(&todos, `SELECT * FROM todo WHERE username = ? AND id IN(
            SELECT rowid FROM todo_search WHERE todo_search MATCH ('category:' || ?)
        )`, username, category)
		fmt.Println(todos)
		return todos, err
	}
	if title != "" {
		err := db.Select(&todos, `SELECT * FROM todo WHERE username = ? AND id IN (
            SELECT rowid FROM todo_search WHERE todo_search MATCH ('title:' || ?)
        )`, username, title)
		fmt.Println(todos)
		return todos, err
	}

	err := db.Select(&todos, "SELECT * FROM todo WHERE username = ?", username)
	return todos, err
}

func DeleteTodo(db *sqlx.DB, id int, username string) error {
	result, err := db.Exec("DELETE FROM todo WHERE username = ? AND id = ?", username, id)
	if err != nil {
		return err
	}

	row_affected_count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row_affected_count == 0 {
		return fmt.Errorf("User %s has no todo with id %d", username, id)
	}

	return nil
}
