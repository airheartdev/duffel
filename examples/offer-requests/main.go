// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/thetreep/duffel"
)

func main() {
	apiToken := os.Getenv("DUFFEL_TOKEN")
	client := duffel.New(apiToken)

	// childAge := 1

	data, err := client.CreateOfferRequest(context.Background(), duffel.OfferRequestInput{
		ReturnOffers: true,

		Passengers: []duffel.OfferRequestPassenger{
			{
				// FamilyName: "Earhart",
				// GivenName:  "Amelia",
				Type: duffel.PassengerTypeAdult,
				// LoyaltyProgrammeAccounts: []duffel.LoyaltyProgrammeAccount{
				// 	{
				// 		AirlineIATACode: "QF",
				// 		AccountNumber:   "1922223336",
				// 	},
				// },
			},
			// {
			// 	Type: duffel.PassengerTypeAdult,
			// },
			// {
			// 	Age: childAge,
			// },
		},
		CabinClass: duffel.CabinClassEconomy,
		Slices: []duffel.OfferRequestSlice{
			{
				DepartureDate: duffel.Date(time.Date(2022, time.June, 24, 0, 0, 0, 0, time.UTC)),
				Origin:        "AUS",
				Destination:   "MSP",
			},
			{
				DepartureDate: duffel.Date(time.Date(2022, time.June, 26, 0, 0, 0, 0, time.UTC)),
				Origin:        "MSP",
				Destination:   "AUS",
			},
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Received %d flight offers:\n\n", len(data.Offers))

	for _, slice := range data.Slices {
		fmt.Printf("\t%s -> %s on %s\n", slice.Origin.Name, slice.Destination.Name, slice.DepartureDate.String())
	}

	fmt.Println()

	for _, offer := range data.Offers {
		fmt.Printf("===> Offer %s from %s\n     Passengers: ", offer.ID, offer.Owner.Name)
		for i, p := range offer.Passengers {
			fmt.Printf("(%s) %s %s", p.Type, p.GivenName, p.FamilyName)
			if i < len(offer.Passengers)-1 {
				fmt.Print(", ")
			}
		}

		fmt.Println()

		fmt.Printf("---> Flights $%s\n", offer.TotalAmount().String())
		for _, s := range offer.Slices {
			fmt.Printf("    🛫 %s to %s\n", *s.Origin.CityName, *s.Destination.CityName)

			for _, f := range s.Segments {
				fmt.Printf("    Departing at %s • Arriving at %s\n", f.DepartingAt, f.ArrivingAt)
			}

			fmt.Printf("    🛬 %s • %s\n", s.FareBrandName, time.Duration(s.Duration).String())
		}
		fmt.Println()
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
