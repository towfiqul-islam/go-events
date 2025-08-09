package models

import (
	"errors"

	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := `INSERT INTO users(email, password) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(u.Email, u.Password)

	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()

	u.ID = userId

	return err
}

func (u *User) ValidateUser() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string

	err := row.Scan(&u.ID, &retrievedPassword)

	if err != nil {
		return err
	}

	isPasswordValid := utils.CheckHashPassword(u.Password, retrievedPassword)

	if !isPasswordValid {
		return errors.New("invalid creds")
	}

	return nil

}

func GetUser(userId int64) (*User, error) {
	query := `SELECT * FROM users WHERE id = ?`
	row := db.DB.QueryRow(query, userId)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
