package main

import (
    "fmt"
	"todo-app/database"
	"todo-app/server"
)

func main() {
	db, err := database.Init()
	if err != nil {
		return
	}
    server := server.NewServer(3000, db)
    err = server.ListenAndServe()
    if err != nil {
        panic(fmt.Sprintf("cannot start server: %s", err))
    }
}
