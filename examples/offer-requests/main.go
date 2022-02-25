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
	client := duffel.New(apiToken)

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
				DepartureDate: duffel.Date(time.Now().AddDate(0, 0, 7)),
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

		log.Printf("Flights $%s", offer.TaxAmount)
		for _, s := range offer.Slices {
			log.Printf("- %s to %s on %s\n", *s.Origin.CityName, *s.Destination.CityName, time.Time(s.DepartureDate).Format("Mon Jan 2 15:04"))

			distances := collect(s.Segments, func(t duffel.Flight) float64 {
				return float64(t.Distance)
			})

			log.Printf("distance: %1.2fkm", sum(distances))
		}
	}
}

func sum[T int | int64 | float64](nums []T) T {
	s := T(0)
	for _, num := range nums {
		s += num
	}
	return s
}

func collect[T any, R any](items []T, f func(T) R) []R {
	out := make([]R, len(items))
	for i, item := range items {
		out[i] = f(item)
	}
	return out
}
