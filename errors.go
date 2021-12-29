package duffel

import "fmt"

type DuffelErr error

var ErrNotFound DuffelErr = fmt.Errorf("duffel: not found")
