package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/airheartdev/duffel"
)

func main() {
	apiToken := os.Getenv("DUFFEL_TOKEN")
	client := duffel.New(apiToken)

	ctx := context.Background()

	if len(os.Args) < 2 {
		log.Println("query is required")
		log.Fatalln("Usage: places <query>")
		return
	}
	query := os.Args[1]

	places, err := client.PlaceSuggestions(ctx, query)
	if err != nil {
		log.Fatalln(err)
		return
	}

	for _, place := range places {
		fmt.Printf("- %s (%s)\n", place.Name, place.IATACode)
	}
}
