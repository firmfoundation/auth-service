package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/firmfoundation/auth-service/cmd/model"

	_ "github.com/jackc/pgx"
	"github.com/jackc/pgx/v4"
)

const webPort = "28280"

type Config struct {
	DB     *pgx.Conn
	Models model.Models
}

func main() {
	log.Println("Starting authentication service..", webPort)

	//connect to db
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to postgresql service")
	}

	app := Config{
		DB:     conn,
		Models: model.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dns string) (*pgx.Conn, error) {
	//db, err := sql.Open("pgx", dns)
	//urlExample := "postgres://postgres:password@127.0.0.1:5432/users"
	conn, err := pgx.Connect(context.Background(), dns)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func connectToDB() *pgx.Conn {
	var counts int
	dns := os.Getenv("DNS")

	for {
		connection, err := openDB(dns)

		if counts > 10 {
			log.Println(err)
			return nil
		}

		if err != nil {
			log.Println("postgresql service not yet ready ...")
			counts++
			log.Println("backing off for two seconds ...")
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Connected to postgres service!")
		return connection
	}
}
