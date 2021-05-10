package main

import (
	"log"
	"os"

	"github.com/danilomarques1/personalfinance/api/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables %v", err)
	}
	var app server.App

	dbCon := server.DbConn{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
	}
	app.Initialize(dbCon)

	app.Start()
}
