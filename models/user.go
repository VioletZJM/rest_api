package models

import (
	"errors"

	"example.com/rest_api/database"
	"example.com/rest_api/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := `
	INSERT INTO users(email,password) VALUES (?,?)
	`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashPassword) //传值给sql

	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()

	u.ID = userId
	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id,password FROM users WHERE email = ?"
	row := database.DB.QueryRow(query, u.Email)

	var retrievedPasswprd string

	err := row.Scan(&u.ID, &retrievedPasswprd)
	if err != nil {
		return errors.New("Credentials invalid")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPasswprd)
	if !passwordIsValid {
		return errors.New("Credentials invalid")
	}

	return nil
}
