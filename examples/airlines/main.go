package main

import (
	"context"
	"log"
	"os"

	"github.com/airheartdev/duffel"
	"github.com/gocarina/gocsv"
)

func main() {
	ctx := context.Background()
	client := duffel.New(os.Getenv("DUFFEL_TOKEN"))

	airlinesCache := map[string]*duffel.Airline{}

	// airline, err := client.GetAirline(ctx, "LV")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Printf("%s (%s) - %s\n", airline.Name, airline.IATACode, airline.ID)

	iter := client.ListAirlines(ctx)

	for iter.Next() {
		airline := iter.Current()
		airlinesCache[airline.ID] = airline
		// fmt.Printf("%s - %s\n", airline.IATACode, airline.Name)
	}
	if iter.Err() != nil {
		log.Fatalln(iter.Err())
	}

	log.Println("Loaded all airlines", len(airlinesCache))

	airlinesCSVFile, err := os.Create("airlines.csv")
	if err != nil {
		log.Fatalln(err)
	}

	airlineRows := []*duffel.Airline{}
	for _, airline := range airlinesCache {
		airlineRows = append(airlineRows, airline)
	}

	err = gocsv.MarshalFile(airlineRows, airlinesCSVFile)
	if err != nil {
		log.Fatalln(err)
	}

	// for _, airline := range airlinesCache {
	// 	if airline.IATACode == "LV" {
	// 		fmt.Printf("%s (%s) - %s\n", airline.ID, airline.Name, airline.IATACode)
	// 	}
	// }
}
