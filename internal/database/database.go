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
    is_done INTEGER NOT NULL
);
CREATE VIRTUAL TABLE IF NOT EXISTS todo_search using fts5 (category, title);

DROP TRIGGER IF EXISTS sync_todo_insert;
DROP TRIGGER IF EXISTS sync_todo_update;
DROP TRIGGER IF EXISTS sync_todo_delete;

CREATE TRIGGER IF NOT EXISTS sync_todo_insert
AFTER INSERT ON todo 
BEGIN 
INSERT OR IGNORE INTO todo_search (category, title) VALUES (NEW.category, NEW.title);
END;

CREATE TRIGGER IF NOT EXISTS sync_todo_update
AFTER UPDATE ON todo 
BEGIN 
UPDATE todo_search SET category = NEW.category, title = NEW.title WHERE rowid = OLD.rowid;
END;

CREATE TRIGGER IF NOT EXISTS sync_todo_delete
AFTER DELETE ON todo 
BEGIN 
DELETE FROM todo_search where rowid = OLD.rowid;
END;`

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

type Todo struct {
	Id          int    `db:"id" json:"id,string"`
	Username    string `db:"username" json:"-"`
	Category    string `db:"category" json:"category"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	IsDone      bool   `db:"is_done" json:"is_done"`
}

func Init() (*sqlx.DB, error) {
	sqlite, err := sqlx.Connect("sqlite3", "todo.db")
	if err != nil {
		return nil, err
	}
	sqlite.MustExec(schema)
	return sqlite, nil
}
