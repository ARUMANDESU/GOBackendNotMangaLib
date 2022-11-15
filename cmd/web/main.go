package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"notmangalib.com/internal/models"
	"os"
	"time"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	manga    *models.MangaModel
	user     *models.UserModel
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dataBasePassword := os.Getenv("database_password")
	dataBaseName := os.Getenv("database_name")
	dbConn, dbErr := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://postgres:%s@localhost:5432/%s", dataBasePassword, dataBaseName))
	if dbErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbErr)
		os.Exit(1)
	}
	defer dbConn.Close()
	var greeting string
	dbErr = dbConn.QueryRow(context.Background(), "select 'DB connected!'").Scan(&greeting)

	if dbErr != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", dbErr)
		os.Exit(1)
	}
	fmt.Println(greeting)

	addr := flag.String("addr", ":5000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// Initialize a new instance of our application struct, containing the
	// dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		manga:    &models.MangaModel{DB: dbConn},
		user:     &models.UserModel{DB: dbConn},
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,

		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,

		MaxHeaderBytes: 524288,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe() // if you want make it https => err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)

}
