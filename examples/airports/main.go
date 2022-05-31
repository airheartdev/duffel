package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/thetreep/duffel"
)

func main() {
	ctx := context.Background()
	client := duffel.New(os.Getenv("DUFFEL_TOKEN"))
	iter := client.ListAirports(ctx, duffel.ListAirportsParams{
		// IATACountryCode: "AU",
	})

	cache := map[string]*duffel.Airport{}

	for iter.Next() {
		airport := iter.Current()
		cache[airport.ID] = airport
		fmt.Printf("%s (%s) - %s, %s\n", airport.Name, airport.IATACode, airport.CityName, airport.IATACountryCode)
	}
	if iter.Err() != nil {
		log.Fatalln(iter.Err())
	}

	log.Println("Loaded all airports", len(cache))

	csvFilePath, err := os.Create("examples/airports/airports.json")
	if err != nil {
		log.Fatalln(err)
	}

	rows := []*duffel.Airport{}
	for _, airline := range cache {
		rows = append(rows, airline)
	}

	err = json.NewEncoder(csvFilePath).Encode(rows)
	if err != nil {
		log.Fatalln(err)
	}
}
