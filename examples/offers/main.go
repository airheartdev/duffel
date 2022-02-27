package main

import (
	"context"
	"log"
	"os"

	"github.com/airheartdev/duffel"
	"github.com/urfave/cli/v2"
)

func main() {
	ctx := context.Background()
	app := &cli.App{
		Name:  "flights",
		Usage: "CLI interface for managing flight bookings",

		Commands: []*cli.Command{
			{
				Name:    "offer-requests",
				Aliases: []string{"or"},
				Usage:   "Manage offer requests",
				Subcommands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "List offer requests",
						Action:  listOfferRequests,
					},
				},
			},
			{
				Name:    "offers",
				Aliases: []string{"of"},
				Usage:   "Manage flight offers",
				Subcommands: []*cli.Command{
					{
						Name:      "list",
						Action:    listOffersAction,
						ArgsUsage: "OFFER_REQUEST_ID",
						Aliases:   []string{"l"},
						Usage:     "List offers for a given request",
					},
				},
			},
		},
	}

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func listOfferRequests(c *cli.Context) error {
	client := duffel.New(os.Getenv("DUFFEL_TOKEN"))

	iter := client.ListOfferRequests(c.Context)

	for iter.Next() {
		req := iter.Current()
		log.Println(req.ID, req.CreatedAt.String())
	}

	return iter.Err()
}

func listOffersAction(c *cli.Context) error {
	client := duffel.New(os.Getenv("DUFFEL_TOKEN"))

	requestID := c.Args().First()
	iter := client.ListOffers(c.Context, requestID, duffel.ListOffersParams{
		MaxConnections: 1,
		Sort:           duffel.ListOffersSortTotalAmount,
	})

	for iter.Next() {
		offer := iter.Current()
		log.Printf("Offer: %s", offer.ID)
		log.Println(offer.ID, offer.Owner.Name)
	}

	return iter.Err()
}
