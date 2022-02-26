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
	iter := client.ListAirports(ctx, duffel.ListAirportsParams{
		IATACountryCode: "AU",
	})

	for iter.Next() {
		airport := iter.Current()
		fmt.Printf("%s\n", airport.Name)
	}

	if iter.Err() != nil {
		log.Fatalln(iter.Err())
	}
}
