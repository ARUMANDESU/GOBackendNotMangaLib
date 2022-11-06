package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

func OpenDB(dbName string, password string) (*pgxpool.Pool, error) {
	dbConn, dbErr := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://postgres:%+v@localhost:5432/%+v", password, dbName))
	if dbErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbErr)
		os.Exit(1)
		return nil, dbErr
	}
	defer dbConn.Close()
	var greeting string
	dbErr = dbConn.QueryRow(context.Background(), "select 'DB connected!'").Scan(&greeting)

	if dbErr != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", dbErr)
		os.Exit(1)
		return nil, dbErr
	}
	fmt.Println(greeting)

	return dbConn, nil
}
