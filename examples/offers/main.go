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

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/thetreep/duffel"
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
					{
						Name:    "create",
						Aliases: []string{"c"},
						Usage:   "Create offer request",
						Action:  createOfferRequest,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "origin",
								Usage:    "Origin airport code",
								Required: true,
								Aliases:  []string{"o"},
							},
							&cli.StringFlag{
								Name:     "destination",
								Usage:    "Destination airport code",
								Required: true,
								Aliases:  []string{"d"},
							},
							&cli.StringFlag{
								Name:     "departure-date",
								Aliases:  []string{"dep"},
								Usage:    "Departure date",
								Required: true,
							},

							&cli.StringFlag{
								Name:     "return-date",
								Aliases:  []string{"ret"},
								Usage:    "Return date",
								Required: false,
							},

							&cli.IntFlag{
								Name:  "adults",
								Usage: "Number of adults",
							},

							&cli.IntSliceFlag{
								Name:    "child-ages",
								Aliases: []string{"child"},
								Usage:   "Age of each child",
							},
						},
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
					{
						Name:      "get",
						Action:    getOfferAction,
						ArgsUsage: "OFFER_ID",
						Usage:     "Get a single offer by ID e.g. off_0000AKzq9VM7JGL9kddOcy",
					},
					{
						Name:      "seats",
						Action:    getOfferSeatsAction,
						ArgsUsage: "OFFER_ID",
						Usage:     "Get seats for single offer by ID e.g. off_0000AKzq9VM7JGL9kddOcy",
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

func createOfferRequest(c *cli.Context) error {
	client := duffel.New(os.Getenv("DUFFEL_TOKEN"))

	slices := []duffel.OfferRequestSlice{}
	origin := c.String("origin")
	destination := c.String("destination")
	adultsCount := c.Int("adults")
	childAges := c.IntSlice("child-ages")

	departureDateStr := c.String("departure-date")
	returnDateStr := c.String("return-date")

	departureDate, err := time.Parse(duffel.DateFormat, departureDateStr)
	if err != nil {
		return err
	}

	slices = append(slices, duffel.OfferRequestSlice{
		Origin:        origin,
		Destination:   destination,
		DepartureDate: duffel.Date(departureDate),
	})

	if returnDateStr != "" {
		returnDate, err := time.Parse(duffel.DateFormat, returnDateStr)
		if err != nil {
			return err
		}
		slices = append(slices, duffel.OfferRequestSlice{
			Origin:        destination,
			Destination:   origin,
			DepartureDate: duffel.Date(returnDate),
		})
	}

	passengers := []duffel.OfferRequestPassenger{}
	for i := 0; i < adultsCount; i++ {
		passengers = append(passengers, duffel.OfferRequestPassenger{
			Type: duffel.PassengerTypeAdult,
		})
	}

	for _, age := range childAges {
		passengers = append(passengers, duffel.OfferRequestPassenger{
			// Type: duffel.PassengerTypeChild,
			Age: age,
		})
	}

	request, err := client.CreateOfferRequest(c.Context, duffel.OfferRequestInput{
		Slices:     slices,
		Passengers: passengers,
	})
	if err != nil {
		return err
	}

	log.Printf("Request ID: %s", request.ID)
	for _, slice := range request.Slices {
		log.Printf("Request slice: %s to %s departing on %s", slice.Origin.Name, slice.Destination.Name, slice.DepartureDate.String())
	}

	return nil
}

func listOfferRequests(c *cli.Context) error {
	client := duffel.New(os.Getenv("DUFFEL_TOKEN"))

	iter := client.ListOfferRequests(c.Context)

	// t := table.NewWriter()
	// t.SetOutputMirror(os.Stdout)
	// t.SetStyle(table.StyleColoredBright)
	// t.AppendHeader(table.Row{
	// 	"ID", "Origin", "Destination", "Departure Date", "Return Date", "Status",
	// })

	for iter.Next() {
		req := iter.Current()
		fmt.Printf("===> Offer Request: %s created: %s\n", req.ID, time.Time(req.CreatedAt).Format(time.RFC3339))

		for _, slice := range req.Slices {
			fmt.Printf("   > %s to %s on %s\n", slice.Origin.IATACode, slice.Destination.IATACode, slice.DepartureDate.String())
		}
	}

	return iter.Err()
}

func getOfferAction(c *cli.Context) error {
	client := duffel.New(os.Getenv("DUFFEL_TOKEN"))
	offerID := c.Args().First()

	off, err := client.GetOffer(c.Context, offerID, duffel.GetOfferParams{
		ReturnAvailableServices: true,
	})
	if err != nil {
		log.Printf("OfferID: %s", offerID)
		return err
	}

	fmt.Printf("Offer: %s\n", offerID)

	fmt.Println("Available services:")
	for _, service := range off.AvailableServices {
		fmt.Printf("  > %s segments: %+v price: %s\n", service.Type, service.SegmentIDs, service.RawTotalAmount)
	}

	fmt.Println("Conditions:")
	for _, slc := range off.Slices {
		fmt.Printf("  > %s - %s %s changeable: %v\n", slc.Origin.Name, slc.Destination.Name, slc.FareBrandName, slc.Changeable)

		if slc.Conditions.ChangeBeforeDeparture != nil {
			fmt.Printf("  > Change before departure: %v\n", slc.Conditions.ChangeBeforeDeparture.Allowed)
		}
	}

	return nil
}

func getOfferSeatsAction(c *cli.Context) error {
	client := duffel.New(os.Getenv("DUFFEL_TOKEN"))
	offerID := c.Args().First()

	seats, err := client.GetSeatmap(c.Context, offerID)
	if err != nil {
		log.Printf("OfferID: %s", offerID)
		return err
	}

	fmt.Printf("Offer: %s\n", offerID)

	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	fmt.Println("Seatmap:")
	for _, seat := range seats {
		fmt.Printf("Slice: %s\n", seat.SliceID)
		for _, cabin := range seat.Cabins {
			fmt.Printf("%s\n", cabin.CabinClass)

			for _, row := range cabin.Rows {
				fmt.Printf("\n")
				for _, section := range row.Sections {
					for _, el := range section.Elements {
						if len(el.AvailableServices) > 0 {
							fmt.Print(green.Sprintf("%-3s ", el.Designator))
						} else {
							fmt.Print(red.Sprintf("%-3s ", el.Designator))
						}
					}
				}
			}
		}

	}

	return nil
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
		fmt.Printf("===> Offer: %s %s\n", offer.ID, offer.Owner.Name)
	}

	return iter.Err()
}
