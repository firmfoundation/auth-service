package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/firmfoundation/auth-service/cmd/model"
	"github.com/joho/godotenv"

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

	//load env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
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

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(url string) (*pgx.Conn, error) {
	//db, err := sql.Open("pgx", dns)
	conn, err := pgx.Connect(context.Background(), url)
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
	url := "postgres://" + os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	for {
		connection, err := openDB(url)

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
