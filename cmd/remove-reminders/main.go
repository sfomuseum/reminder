package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sfomuseum/reminder/app/reminders/remove"
)

func main() {

	ctx := context.Background()
	err := remove.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to remove reminders, %v", err)
	}
}
