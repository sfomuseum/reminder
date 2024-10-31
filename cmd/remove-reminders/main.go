package main

import (
	"context"
	"log"

	"github.com/sfomuseum/reminder/app/reminders/remove"
)

func main() {

	ctx := context.Background()
	err := remove.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to remove reminders, %v", err)
	}
}
