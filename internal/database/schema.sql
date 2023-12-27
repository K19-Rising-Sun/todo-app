CREATE TABLE IF NOT EXISTS users (
    username TEXT NOT NULL PRIMARY KEY, 
    password TEXT NOT NULL
 );
CREATE TABLE IF NOT EXISTS todos (
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
AFTER INSERT ON todos 
BEGIN 
INSERT OR IGNORE INTO todo_search (category, title) VALUES (NEW.category, NEW.title);
END;

CREATE TRIGGER IF NOT EXISTS sync_todo_update
AFTER UPDATE ON todos 
BEGIN 
UPDATE todo_search SET category = NEW.category, title = NEW.title WHERE rowid = OLD.rowid;
END;

CREATE TRIGGER IF NOT EXISTS sync_todo_delete
AFTER DELETE ON todos 
BEGIN 
DELETE FROM todo_search where rowid = OLD.rowid;
END;
