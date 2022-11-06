package models

import "github.com/jackc/pgx/v4/pgxpool"

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(name string, email string, hashedPassword string) (*User, error) {
	user := &User{1, "sdf", "sdf"}
	return user, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	return nil, nil
}
