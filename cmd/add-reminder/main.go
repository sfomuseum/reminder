package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sfomuseum/reminder/app/reminders/add"
)

func main() {

	ctx := context.Background()
	err := add.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to add reminder, %v", err)
	}
}
