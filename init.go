package main

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Mongo struct {
	}
	Cassandra struct{}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectMongo()
}
