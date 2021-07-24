package database

import (
	"context"
	"fmt"
	"instagram/config"
	"instagram/ent"
	"log"

	_ "github.com/lib/pq"
)

var Client *ent.Client

func CreateConnection() {

	client, err := ent.Open(config.Get("DB_DRIVER"), fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", config.Get("DB_HOST"), config.Get("DB_PORT"), config.Get("DB_USER"), config.Get("DB_NAME"), config.Get("DB_PASS")))

	if err != nil {
		log.Fatalf("failed opening connection to db: %v", err)
	}
	// apply migrations
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	Client = client
	fmt.Println("connected")
}
