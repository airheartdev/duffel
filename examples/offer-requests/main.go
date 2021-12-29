package main

import (
	"context"
	"log"

	"github.com/airheartdev/duffel"
)

func main() {
	client := duffel.New("token", duffel.WithDefaultAPI())
	offers, err := client.OfferRequest(context.Background(), &duffel.OfferRequestInput{
		Passengers: []duffel.Passenger{
			{
				FamilyName: "Earhart",
				GivenName:  "Amelia",
				Type:       duffel.PassengerTypeAdult,
			},
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(offers)
}
