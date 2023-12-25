package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var schema = `
CREATE TABLE IF NOT EXISTS user (
    username TEXT NOT NULL PRIMARY KEY, 
    password TEXT NOT NULL
 );
CREATE TABLE IF NOT EXISTS todo (
    id INTEGER PRIMARY KEY, 
    username TEXT NOT NULL, 
    category TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    state TEXT NOT NULL
);
CREATE VIRTUAL TABLE IF NOT EXISTS todo_search using fts5 (category, title);

CREATE TRIGGER IF NOT EXISTS insert_todo_trigger
AFTER INSERT ON todo 
BEGIN 
INSERT OR IGNORE INTO todo_search (category, title) VALUES (NEW.category, NEW.title);
END;

CREATE TRIGGER IF NOT EXISTS update_todo_trigger
AFTER UPDATE ON todo 
BEGIN 
UPDATE todo_search SET category = NEW.TITLE, title = NEW.TITLE WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS delete_todo_trigger
AFTER DELETE ON todo 
BEGIN 
DELETE FROM todo_search where rowid = OLD.rowid;
END;`

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
