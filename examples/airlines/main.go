package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/airheartdev/duffel"
)

func main() {
	ctx := context.Background()
	client := duffel.New(os.Getenv("DUFFEL_TOKEN"))
	iter := client.ListAirlines(ctx)

	for iter.Next() {
		airline := iter.Current()
		fmt.Printf("%s - %s\n", airline.IATACode, airline.Name)
	}

	if iter.Err() != nil {
		log.Fatalln(iter.Err())
	}
}
