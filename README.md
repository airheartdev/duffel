# Duffel API Go Client

A Go (golang) client library for the [Duffel](https://duffel.com) API implemented by the Airheart team.

[![Tests](https://github.com/airheartdev/duffel/actions/workflows/ci.yaml/badge.svg)](https://github.com/airheartdev/duffel/actions/workflows/ci.yaml)

## Installation

**Requires at least Go 1.18-rc1 since we use generics on the internal API client**

```shell
go get github.com/airheartdev/duffel
```

## Usage examples

See the [examples/\*](/examples/) directory

## Usage

### func \(\*API\) CreateOfferRequest

```go
func (a *API) CreateOfferRequest(ctx context.Context, requestInput OfferRequestInput) (*OfferRequest, error)
```

### func \(\*API\) CreateOrder

```go
func (a *API) CreateOrder(ctx context.Context, input CreateOrderInput) (*Order, error)
```

CreateOrder creates a new order\.

### func \(\*API\) GetAircraft

```go
func (a *API) GetAircraft(ctx context.Context, id string) (*Aircraft, error)
```

### func \(\*API\) GetAirline

```go
func (a *API) GetAirline(ctx context.Context, id string) (*Airline, error)
```

### func \(\*API\) GetAirport

```go
func (a *API) GetAirport(ctx context.Context, id string) (*Airport, error)
```

### func \(\*API\) GetOffer

```go
func (a *API) GetOffer(ctx context.Context, id string) (*Offer, error)
```

GetOffer gets a single offer by ID\.

### func \(\*API\) GetOfferRequest

```go
func (a *API) GetOfferRequest(ctx context.Context, id string) (*OfferRequest, error)
```

### func \(\*API\) GetOrder

```go
func (a *API) GetOrder(ctx context.Context, id string) (*Order, error)
```

CreateOrder creates a new order\.

### func \(\*API\) GetSeatmaps

```go
func (a *API) GetSeatmaps(ctx context.Context, offerID string) *Iter[Seatmap]
```

### func \(\*API\) ListAircraft

```go
func (a *API) ListAircraft(ctx context.Context, params ...ListAirportsParams) *Iter[Aircraft]
```

### func \(\*API\) ListAirlines

```go
func (a *API) ListAirlines(ctx context.Context, params ...ListAirportsParams) *Iter[Airline]
```

### func \(\*API\) ListAirports

```go
func (a *API) ListAirports(ctx context.Context, params ...ListAirportsParams) *Iter[Airport]
```

### func \(\*API\) ListOfferRequests

```go
func (a *API) ListOfferRequests(ctx context.Context) *Iter[OfferRequest]
```

### func \(\*API\) ListOffers

```go
func (a *API) ListOffers(ctx context.Context, offerRequestId string, options ...ListOffersParams) *Iter[Offer]
```

ListOffers lists all the offers for an offer request\. Returns an iterator\.

### func \(\*API\) ListOrders

```go
func (a *API) ListOrders(ctx context.Context, params ...ListOrdersParams) *Iter[Order]
```

### func \(\*API\) UpdateOfferPassenger

```go
func (a *API) UpdateOfferPassenger(ctx context.Context, offerRequestID, passengerID string, input *PassengerUpdateInput) (*OfferRequestPassenger, error)
```

UpdateOfferPassenger updates a single offer passenger\.

### func \(\*API\) UpdateOrder

```go
func (a *API) UpdateOrder(ctx context.Context, id string, params OrderUpdateParams) (*Order, error)
```

### Monetary amounts

All models that have fields for an amount and a currency implement a method to return a parsed `currency.Amount{}`

```go
total:= order.TotalAmount() // currency.Amount
total.Currency // USD
total.Amount // 100.00
total.String() // 100.00 USD
```

### Working with iterators

All requests that return more than one record will return an iterator. An Iterator automatically paginates results and respects rate limits, reducing the complexity of the overall programming model.

```go
client := duffel.New(os.Getenv("DUFFEL_TOKEN"))
iter := client.ListAirports(ctx, duffel.ListAirportsParams{
  IATACountryCode: "AU",
})

for iter.Next() {
  airport := iter.Current()
  fmt.Printf("%s\n", airport.Name)
}

// If there was an error, the loop above will exit so that you can
// inspect the error here:
if iter.Err() != nil {
  log.Fatalln(iter.Err())
}
```

We've added a convenience method to collect all items from the iterator into a single slice in one go:

```go
airports, err := duffel.Collect(iter)
if err != nil  {
  // Handle error from iter.Err()
}

// airports is a []*duffel.Airport
```

### Error Handling

Each API method returns an error or an iterator that returns errors at each iteration. If an error is returned from Duffel, it will be of type `DuffelError` and expose more details on how to handle it.

```go
// Example error inspection after making an API call:
offer, err := client.GetOffer(ctx, "off_123")
if err != nil {
  if err, ok:= err.(duffel.DuffelError); ok {
    // err.Errors[0].Type etc
    // err.IsCode(duffel.BadRequest)
  }else{
    // Do something with regular Go error
  }
}
```

## Implementation status:

To maintain simplicity and ease of use, this client library is hand-coded (instead of using Postman to Go code generation) and contributions are greatly apprecicated.

- [x] Most API types
- [x] API Client
- [x] Error handling
- [x] Pagination _(using iterators)_
- [x] Rate Limiting _(automatically throttles requests to stay under limit)_
- [x] Offer Requests
  - [x] Create offer request and return offer
  - [x] Get offer by ID
  - [x] List all offers
- [x] Offers
- [x] Orders
- [ ] Payments
- [x] Seat Maps
- [ ] Order Cancellations
- [ ] Order Change Requests
- [ ] Order Change Offers
- [ ] Order Changes
- [x] Airports
- [x] Airlines
- [x] Equipment (Aircraft)

## License

MIT.
