package main

import (
	"context"
	"log"

	"github.com/sfomuseum/reminder/app/reminders/add"
)

func main() {

	ctx := context.Background()
	err := add.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to add reminder, %v", err)
	}
}
