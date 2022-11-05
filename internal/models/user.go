package models

import "github.com/jackc/pgx/v4/pgxpool"

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(title string, description string, author string, mangaType string) (*User, error) {
	user := &User{1, "sdf"}
	return user, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	return nil, nil
}
