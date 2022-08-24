package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/thetreep/duffel"
)

func main() {
	ctx := context.Background()
	client := duffel.New(os.Getenv("DUFFEL_TOKEN"))

	iter := client.Cities(ctx)

	count := 0

	for iter.Next() {
		count++
		city := iter.Current()
		fmt.Printf("%s (%s) - %s, %s\n", city.ID, city.Name, city.IATACode, *city.IATACountryCode)
	}

	fmt.Printf("Total cities: %d", count)

	if iter.Err() != nil {
		log.Fatalln(iter.Err())
	}
}
