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

_TBD_

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
