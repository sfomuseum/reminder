package main

import (
	"context"
	"log"

	"github.com/sfomuseum/reminder/app/dynamodb/tables/create"
)

func main() {

	ctx := context.Background()
	err := create.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to create tables, %v", err)
	}
}
