package service

import (
	"todo-app/database"

	"github.com/jmoiron/sqlx"
)

func AddUser(db *sqlx.DB, user *database.User) error {
	_, err := db.NamedExec("INSERT INTO user (username, password) VALUES (:username, :password)", user)
	return err
}

func IsValidUser(db *sqlx.DB, user *database.User) (bool, error) {
	users := []database.User{}
	err := db.Select(&users, "SELECT * FROM user WHERE username = $1", user.Username)
	if err != nil {
		return false, err
	}
	if len(users) != 1 {
		return false, nil
	}
    return users[0].Password == user.Password, nil
}
