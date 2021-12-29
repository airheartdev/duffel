package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/airheartdev/duffel"
)

func main() {
	apiToken := os.Getenv("DUFFEL_TOKEN")
	client := duffel.New(apiToken, duffel.WithDefaultAPI())

	adult := duffel.PassengerTypeAdult

	offers, err := client.OfferRequest(context.Background(), &duffel.OfferRequestInput{
		ReturnOffers: true,
		Passengers: []duffel.Passenger{
			{
				FamilyName: "Earhart",
				GivenName:  "Amelia",
				Type:       &adult,
			},
		},
		CabinClass: duffel.CabinClassEconomy,
		Slices: []duffel.OfferRequestSlice{
			{
				DepartureDate: duffel.Date(time.Now().AddDate(0, 0, 1)),
				Origin:        "JFK",
				Destination:   "AUS",
			},
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(offers)
}
