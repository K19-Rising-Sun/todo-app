package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"todo-app/internal/server"
)

//go:embed internal/database/schema.sql
var schema string

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "todo.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := db.ExecContext(ctx, schema); err != nil {
		fmt.Println(err)
		return
	}

	server := server.NewServer(8000, db)
	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
