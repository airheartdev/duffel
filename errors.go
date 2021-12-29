package duffel

import "fmt"

type DuffelError struct {
	Meta   Meta    `json:"meta"`
	Errors []Error `json:"errors"`
}

type Error struct {
	Type             string `json:"type"`
	Title            string `json:"title"`
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
	Code             string `json:"code"`
}

type Meta struct {
	Status    int64  `json:"status"`
	RequestID string `json:"request_id"`
}

type DuffelErr error

var ErrNotFound DuffelErr = fmt.Errorf("duffel: not found")

func buildError(e Error) InvalidRequestErr {
	return InvalidRequestErr(fmt.Errorf("duffel: %s", e.Message))
}

type InvalidRequestErr DuffelErr
