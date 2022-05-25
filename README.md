# Duffel API Go Client

A Go (golang) client library for the [Duffel Flights API](https://duffel.com) implemented by the [Airheart](https://airheart.com) team.

[![Tests](https://github.com/airheartdev/duffel/actions/workflows/ci.yaml/badge.svg)](https://github.com/airheartdev/duffel/actions/workflows/ci.yaml)

## Installation

We've designed this pkg to be familiar and ideomatic to any Go developer. Go get it the usual way:

> NOTE: Requires at least Go 1.18 since we use generics on the internal API client

```shell
go get github.com/airheartdev/duffel
```

## Usage examples

See the [examples/\*](/examples/) directory

## Getting Started

The easiest way to get started, assuming you have a Duffel account set up and an API key ready, is to create an API instance and make your first request:

```go
// Create a new client:
dfl := duffel.New(os.Getenv("DUFFEL_TOKEN"))
```

For available methods, see:

- [GoDoc Documentation](https://pkg.go.dev/github.com/airheartdev/duffel#section-documentation)
- [Duffel API reference](https://duffel.com/docs/api/overview/welcome)

And familiarise yourself with the implementation notes:

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
dfl := duffel.New(os.Getenv("DUFFEL_TOKEN"))
iter := dfl.ListAirports(ctx, duffel.ListAirportsParams{
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

## Error Handling

Each API method returns an error or an iterator that returns errors at each iteration. If an error is returned from Duffel, it will be of type `DuffelError` and expose more details on how to handle it.

```go
// Example error inspection after making an API call:
offer, err := dfl.GetOffer(ctx, "off_123")
if err != nil {
  // Simple error code check
  if duffel.IsErrorCode(err, duffel.AirlineInternal) {
    // Don't retry airline errors, contact support
  }

  // Access the DuffelError object to see more detail
  if derr, ok:= err.(*duffel.DuffelError); ok {
    // derr.Errors[0].Type etc
    // derr.IsCode(duffel.BadRequest)
  }else{
    // Do something with regular Go error
  }
}
```

### `duffel.IsErrorCode(err, code)`

`IsErrorCode` is a concenience method to check if an error is a specific error code from Duffel.
This simplifies error handling branches without needing to type cast multiple times in your code.

### `duffel.IsErrorType(err, typ)`

`IsErrorType` is a concenience method to check if an error is a specific error type from Duffel.
This simplifies error handling branches without needing to type cast multiple times in your code.

You can also check the `derr.Retryable` field, which will be false if you need to contact Duffel support to resolve the issue, and should not be retried. Example, creating an order.

## Implementation status

To maintain simplicity and ease of use, this client library is hand-coded (instead of using Postman to Go code generation) and contributions are greatly apprecicated.

- [x] Types for all API models
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
- [x] Seat Maps
- [x] Order Cancellations
- [x] Order Change Requests
- [x] Order Change Offers
- [x] Order Changes
- [x] Airports
- [x] Airlines
- [x] Equipment (Aircraft)
- [ ] Payments (Looking for contributions)

## License

MIT License
