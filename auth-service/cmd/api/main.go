package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/OlegChuev/microservices/auth/data"
	"github.com/OlegChuev/microservices/utils"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const WEB_PORT = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
	*utils.Config
}

func main() {
	log.Println("Starting auth-service")

	// Connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Cannot connect to PSQL")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)

		if err != nil {
			log.Println("PSQL is not ready")
			counts++
		} else {
			log.Println("Connected to PSQL")
			return connection
		}

		if counts > 10 {
			log.Println(err)

			return nil
		}

		log.Printf("Waiting 2 seconds...")
		time.Sleep(time.Second * 2)
		continue
	}
}
