package main

import (
	"context"
	"log"

	"github.com/sfomuseum/reminder/app/reminders/process"
)

func main() {

	ctx := context.Background()
	err := process.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to process reminders, %v", err)
	}
}
