package models

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	HashedPassword string `json:"hashedPassword"`
}
type Claims struct {
	Id    int
	Name  string
	Email string
	Role  string
	jwt.StandardClaims
}

type SignModel struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (u *UserModel) Insert(name string, email string, hashedPassword string) (*User, error) {
	stmt := `insert into useri(name,email,hashed_password,role) values($1,$2,$3,'user') returning userid, name,email,role ;`
	user := &User{}
	result := u.DB.QueryRow(context.Background(), stmt, name, email, hashedPassword).Scan(&user.Id, &user.Name, &user.Email, &user.Role)
	if result != nil {
		return nil, result
	}

	return user, nil
}

func (u *UserModel) Get(id int) (*User, error) {
	return nil, nil
}

func (u *UserModel) FindCheckUser(email string, password string) (*User, error) {
	stmt := `select * from useri where email=$1;`
	user := &User{}
	result := u.DB.QueryRow(context.Background(), stmt, email).Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword, &user.Role)
	if result != nil {
		return nil, result
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}
