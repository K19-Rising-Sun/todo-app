package database

import (
	"fmt"

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
	fmt.Println(todos)
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
