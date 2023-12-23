package database

import (
	"github.com/jmoiron/sqlx"
)

func (user *User) Add(db *sqlx.DB) error {
	_, err := db.NamedExec("INSERT INTO user (username, password) VALUES (:username, :password)", user)
	return err
}

func (user *User) IsValid(db *sqlx.DB) (bool, error) {
	users := []User{}
	err := db.Select(&users, "SELECT * FROM user WHERE username = $1", user.Username)
	if err != nil {
		return false, err
	}
	if len(users) != 1 {
		return false, nil
	}
	return users[0].Password == user.Password, nil
}
