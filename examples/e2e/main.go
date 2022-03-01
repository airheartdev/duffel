package main

import (
	"context"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/airheartdev/duffel"
)

func main() {
	token := os.Getenv("DUFFEL_TOKEN")

	if token == "" {
		log.Fatalln("DUFFEL_TOKEN is not set")
	} else if !strings.HasPrefix(token, "duffel_test_") {
		log.Fatalln("E2E test cannot be run with a live token")
	}

	ctx := context.Background()

	// Create a new API client
	client := duffel.New(token)

	data, err := client.CreateOfferRequest(ctx, duffel.OfferRequestInput{
		Passengers: []duffel.OfferRequestPassenger{
			{
				FamilyName: "Earhardt",
				GivenName:  "Amelia",
				Type:       duffel.PassengerTypeAdult,
			},
			{
				Age: 14,
			},
		},
		CabinClass:   duffel.CabinClassEconomy,
		ReturnOffers: true,
		Slices: []duffel.OfferRequestSlice{
			{
				DepartureDate: duffel.Date(time.Now().AddDate(0, 0, 7)),
				Origin:        "JFK",
				Destination:   "AUS",
			},
		},
	})
	handleErr(err)

	sort.Slice(data.Offers, func(i, j int) bool {
		cmp, _ := data.Offers[i].TotalAmount().Cmp(data.Offers[j].TotalAmount())
		return cmp < 0
	})

	offer := data.Offers[0]

	log.Println(offer.TotalAmount())

}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
