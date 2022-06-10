// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"

	"github.com/bojanz/currency"
)

type (
	Seatmap struct {
		Cabins    []Cabin `json:"cabins"`
		ID        string  `json:"id"`
		SegmentID string  `json:"segment_id"`
		SliceID   string  `json:"slice_id"`
	}

	Cabin struct {
		Aisles     int        `json:"aisles"`
		CabinClass CabinClass `json:"cabin_class"`
		Deck       int        `json:"deck"`
		// A list of rows in this cabin.
		Rows []Row `json:"rows"`

		// Where the wings of the aircraft are in relation to rows in the cabin.
		Wings Wing `json:"wings"`
	}

	// Row represents a row in a cabin.
	Row struct {
		// A list of sections. Each row is divided into sections by one or more aisles.
		Sections []SeatSection `json:"sections"`
	}

	// SeatSection represents a section of a row.
	SeatSection struct {
		// The elements that make up this section.
		Elements []SectionElement `json:"elements"`
	}

	// SectionElement represents an element in a section.
	SectionElement struct {
		// The element type, e.g. seat, exit_row, stairs, etc.
		Type ElementType `json:"type"`

		// Seats are considered a special kind of service.
		// There will be at most one service per seat per passenger.
		// A seat can only be booked for one passenger. If a seat has no available services (which will be represented as an empty list : []) then it's unavailable.
		AvailableServices []SectionService `json:"available_services"`

		// The designator used to uniquely identify the seat, usually made up of a row number and a column letter
		Designator string `json:"designator"`

		// Each disclosure is text, in English, provided by the airline that describes the terms and conditions of this seat. We recommend showing this in your user interface to make sure that customers understand any restrictions and limitations.
		Disclosures []string `json:"disclosures"`

		// A name which describes the type of seat, which you can display in your user interface to help customers to understand its features.
		// Example: "Exit row seat"
		Name string `json:"name"`
	}

	SectionService struct {
		ID               string `json:"id"`
		PassengerID      string `json:"passenger_id"`
		RawTotalAmount   string `json:"total_amount"`
		RawTotalCurrency string `json:"total_currency"`
	}

	ElementType string

	// Wing represents a wing of the aircraft in relation to rows in the cabin.
	Wing struct {
		// The index of the first row which is overwing, starting from the front of the aircraft.
		FirstRowIndex int `json:"first_row_index"`
		// The index of the last row which is overwing, starting from the front of the aircraft.
		LastRowIndex int `json:"last_row_index"`
	}

	SeatmapClient interface {
		// GetSeatmaps returns an iterator for the seatmaps of a given Offer.
		GetSeatmaps(ctx context.Context, offerID string) *Iter[Seatmap]
	}
)

const (
	ElementTypeSeat     ElementType = "seat"
	ElementTypeBassinet ElementType = "bassinet"
	ElementTypeEmpty    ElementType = "empty"
	ElementTypeExitRow  ElementType = "exit_row"
	ElementTypeLavatory ElementType = "lavatory"
	ElementTypeGalley   ElementType = "galley"
	ElementTypeCloset   ElementType = "closet"
	ElementTypeStairs   ElementType = "stairs"
)

func (e ElementType) String() string {
	return string(e)
}

func (a *API) GetSeatmaps(ctx context.Context, offerID string) *Iter[Seatmap] {
	return newRequestWithAPI[EmptyPayload, Seatmap](a).
		Get("/air/seat_maps").
		WithParam("offer_id", offerID).
		All(ctx)
}

func (s *SectionService) TotalAmount() currency.Amount {
	amount, err := currency.NewAmount(s.RawTotalAmount, s.RawTotalCurrency)
	if err != nil {
		return currency.Amount{}
	}
	return amount
}
