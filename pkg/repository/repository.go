package repository

import "github.com/jackc/pgx/v4/pgxpool"

type Repository struct {
}

type Authorization interface {
	InsertUser(name string, email string, password string, role string)
}

func InitRepository(db *pgxpool.Pool) *Repository {
	return &Repository{}
}