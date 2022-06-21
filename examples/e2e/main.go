// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/airheartdev/duffel"
	"github.com/jedib0t/go-pretty/v6/table"
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
			// {
			// 	Age: 14,
			// },
		},
		CabinClass:   duffel.CabinClassFirst,
		ReturnOffers: false,
		Slices: []duffel.OfferRequestSlice{
			{
				DepartureDate: duffel.Date(time.Now().AddDate(0, 1, 7)),
				Origin:        "JFK",
				Destination:   "AUS",
			},
		},
	})
	handleErr(err)

	offersIter := client.ListOffers(ctx, data.ID)
	var selectedOffer *duffel.Offer

	allOffers, err := duffel.Collect(offersIter)
	handleErr(err)

	for _, offer := range allOffers {
		if offer.Owner.IATACode == "ZZ" {
			selectedOffer = offer
			break
		}
	}

	for i, offer := range allOffers {
		if i > 10 {
			break
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Offer ID", "Total Amount", "Airline", "Flight", "Origin", "Destination", "Cabin Class"})
		t.AppendRow(table.Row{offer.ID, offer.Owner.Name, offer.TotalAmount().String()})
		t.AppendRow(table.Row{"Slice ID", "Origin", "Destination", "Departure Time", "Arrival Time", "Cabin Class", "Changeable"})

		for _, slice := range offer.Slices {
			for _, segment := range slice.Segments {
				dep, _ := segment.DepartingAt()
				arr, _ := segment.ArrivingAt()
				t.AppendRow(table.Row{
					slice.ID,
					segment.Origin.IATACode,
					segment.Destination.IATACode,

					dep.Format(time.RFC822),
					arr.Format(time.RFC822),
					segment.Passengers[0].CabinClass.String(),
					renderChangeableStatus(slice),
				})
			}
		}
		t.AppendSeparator()
		t.SetStyle(table.StyleColoredBright)
		t.Render()
	}
	if err := offersIter.Err(); err != nil {
		handleErr(err)
	}

	if selectedOffer == nil {
		log.Fatalln("No changeable offers found")
	}

	log.Printf("Selected offer: %s", selectedOffer.ID)

	updatedPassengers := []duffel.OrderPassenger{}
	for _, p := range selectedOffer.Passengers {
		pax := duffel.OrderPassenger{
			ID:         p.ID,
			Title:      duffel.PassengerTitleMiss,
			FamilyName: p.FamilyName,
			GivenName:  p.GivenName,
			BornOn:     duffel.Date(time.Now().AddDate(-31, 1, 0)), // 31 years ago
			Email:      "amelia@airheart.com",
			Gender:     duffel.GenderFemale,
			// 						+442080160509
			// PhoneNumber: "+61432221634",
			PhoneNumber: "+16465043216",
			LoyaltyProgrammeAccounts: []duffel.LoyaltyProgrammeAccount{
				{
					AirlineIATACode: "AA",
					AccountNumber:   "AA12345",
				},
			},
		}

		if selectedOffer.PassengerIdentityDocumentsRequired {
			pax.IdentityDocuments = []duffel.IdentityDocument{
				{
					UniqueIdentifier:   "KP1234567",
					ExpiresOn:          duffel.Date(time.Now().AddDate(4, 2, 8)),
					IssuingCountryCode: "GB",
					Type:               "passport",
				},
			}
		}

		updatedPassengers = append(updatedPassengers, pax)
	}

	tableRows := []table.Row{{"Passenger ID", "Family Name", "Given Name", "DoB"}}
	for _, pax := range updatedPassengers {
		tableRows = append(tableRows, []any{pax.ID, pax.FamilyName, pax.GivenName, time.Time(pax.BornOn).Format("2006-01-02")})
	}
	renderTable(tableRows)

	order, err := client.CreateOrder(ctx, duffel.CreateOrderInput{
		Type:           duffel.OrderTypeInstant,
		SelectedOffers: []string{selectedOffer.ID},
		Payments: []duffel.PaymentCreateInput{
			{
				Amount:   selectedOffer.RawTotalAmount,
				Currency: selectedOffer.RawTotalCurrency,
				Type:     duffel.PaymentMethodBalance,
			},
		},
		Metadata: duffel.Metadata{
			"seat_preference": "window",
			"meal_preference": "VGML",
		},
		Passengers: updatedPassengers,
	})

	handleErr(err)

	renderTable([]table.Row{
		{"Order ID", "Reference", "Status", "Total Amount", "Tax Amount", "Airline"},
		{order.ID, order.BookingReference, renderPaymentStatus(order.PaymentStatus), order.TotalAmount().String(), order.TaxAmount().String(), order.Owner.Name},
	})

	sliceRows := []table.Row{{"Origin", "Destination", "Flight Number", "Conditions"}}
	for _, slice := range order.Slices {
		sliceRows = append(sliceRows, table.Row{
			slice.Origin.IATACode,
			slice.Destination.IATACode,
			slice.Segments[0].OperatingCarrierFlightNumber,
			renderChangeableStatus(slice),
		})
	}
	renderTable(sliceRows)

	orderChangeRequest, err := client.CreateOrderChangeRequest(ctx, duffel.OrderChangeRequestParams{
		OrderID: order.ID,
		Slices: duffel.SliceChange{
			Remove: []duffel.SliceRemove{},
			Add: []duffel.SliceAdd{
				{
					DepartureDate: duffel.Date(time.Now().AddDate(0, 2, 7)),
					Origin:        "AUS",
					Destination:   "JFK",
					CabinClass:    duffel.CabinClassFirst,
				},
			},
		},
	})
	handleErr(err)

	log.Printf("Created order change request: %s for order: %s", orderChangeRequest.ID, orderChangeRequest.OrderID)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Offer ID", "Expires", "Total Change Amount", "Penalty Amount"})

	for _, offer := range orderChangeRequest.OrderChangeOffers {
		t.AppendRow(table.Row{offer.ID, offer.ExpiresAt.Format(time.RFC822), offer.ChangeTotalAmount().String(), offer.PenaltyTotalAmount().String()})

		t.AppendRow(table.Row{"Slice ID", "Origin", "Destination", "Departure Time", "Arrival Time", "Cabin Class", "Changeable"})

		for _, slice := range offer.Slices.Add {
			for _, segment := range slice.Segments {
				dep, _ := segment.DepartingAt()
				arr, _ := segment.ArrivingAt()

				t.AppendRow(table.Row{
					fmt.Sprintf("Add %s", slice.ID),
					segment.Origin.IATACode,
					segment.Destination.IATACode,
					dep.Format(time.RFC822),
					arr.Format(time.RFC822),
					// segment.Passengers != nil && segment.Passengers[0].CabinClass.String(),
					renderChangeableStatus(slice),
				})
			}
		}
		for _, slice := range offer.Slices.Remove {
			for _, segment := range slice.Segments {
				dep, _ := segment.DepartingAt()
				arr, _ := segment.ArrivingAt()
				t.AppendRow(table.Row{
					fmt.Sprintf("Remove %s", slice.ID),
					segment.Origin.IATACode,
					segment.Destination.IATACode,
					dep.Format(time.RFC822),
					arr.Format(time.RFC822),
					segment.Passengers[0].CabinClass.String(),
					renderChangeableStatus(slice),
				})
			}
		}
		t.AppendSeparator()
	}

	t.SetStyle(table.StyleColoredBright)
	t.Render()

	orderChange, err := client.CreatePendingOrderChange(ctx, orderChangeRequest.ID)
	handleErr(err)

	log.Printf("Pending order change: %s for order: %s refund to: %s", orderChange.ID, orderChange.OrderID, orderChange.RefundTo.String())

	confirmed, err := client.ConfirmOrderChange(ctx, orderChange.ID, duffel.PaymentCreateInput{
		Amount:   orderChange.RawChangeTotalAmount,
		Currency: orderChange.RawChangeTotalCurrency,
		Type:     duffel.PaymentMethodBalance,
	})
	handleErr(err)

	log.Printf("Confirmed order change: %s for order: %s", confirmed.ID, confirmed.OrderID)
}

func handleErr(err error) {
	if err != nil {
		if err, ok := err.(*duffel.DuffelError); ok {
			log.Fatalf("Duffel API error: %s - %s", err.Errors[0].Code, err.Error())
		} else {
			log.Fatalln(err)
		}
	}
}

func renderTable(rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(rows[0])
	t.AppendRows(rows[1:])
	t.SetStyle(table.StyleColoredBright)
	t.Render()
}

func renderPaymentStatus(p duffel.PaymentStatus) string {
	if p.AwaitingPayment {
		return "awaiting payment"
	}
	return "paid"
}

func renderChangeableStatus(slice duffel.Slice) string {
	if slice.Changeable {
		if slice.Conditions.ChangeBeforeDeparture != nil && slice.Conditions.ChangeBeforeDeparture.Allowed {
			return fmt.Sprintf("Changeable before departure with penalty: %s", slice.Conditions.ChangeBeforeDeparture.PenaltyAmount().String())
		} else {
			return "not changeable before departure"
		}
	}
	return "not changeable"
}

func filter[T any](slice []T, f func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
