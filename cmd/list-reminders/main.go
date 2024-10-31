package main

import (
	"context"
	"log"

	"github.com/sfomuseum/reminder/app/reminders/list"
)

func main() {

	ctx := context.Background()
	err := list.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to remove reminders, %v", err)
	}
}
