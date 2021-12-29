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

	for _, offer := range offers.Offers {
		log.Println("Passengers:")
		for _, p := range offer.Passengers {
			log.Printf("- %s %s (%s)\n", p.GivenName, p.FamilyName, p.Type)
		}

		log.Printf("Flights $%1.2f", offer.TaxAmount)
		for _, s := range offer.Slices {
			log.Printf("- %s to %s on %s\n", *s.Origin.CityName, *s.Destination.CityName, time.Time(s.DepartureDate).Format("Mon Jan 2 15:04"))
		}
	}
}
