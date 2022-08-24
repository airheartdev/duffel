// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

// ListContainer is a general interface for which all list object structs
// should comply. They achieve this by embedding a ListMeta struct and
// inheriting its implementation of this interface.
type ListContainer[T any] interface {
	SetListMeta(*ListMeta)
	GetListMeta() *ListMeta
	GetItems() []*T
	SetItems(items []*T)
	LastRequestID() (string, bool)
}

type List[T any] struct {
	items []*T
	// *Iter[T]
	*ListMeta

	// Duffel Request ID
	lastRequestID string `json:"-" url:"-"`
}

func (l *List[T]) GetItems() []*T {
	return l.items
}

func (l *List[T]) SetItems(items []*T) {
	l.items = items
}

func (l *List[T]) setRequestID(id string) {
	l.lastRequestID = id
}

func (l *List[T]) LastRequestID() (string, bool) {
	return l.lastRequestID, l.lastRequestID != ""
}

func (l *List[T]) SetListMeta(meta *ListMeta) {
	l.ListMeta = meta
}

// ListMeta is the structure that contains the common properties
// of List iterators.
type ListMeta struct {
	// HasMore is a boolean that indicates whether there are more items

	// After is a string that contains the token for the next page of results
	After string `json:"after,omitempty" url:"after,omitempty"`

	// Before is a string that contains the token for the previous page of results
	Before string `json:"before,omitempty" url:"-"`

	// Limit is a number that indicates the maximum number of items to return
	Limit int `json:"limit,omitempty" url:"limit,omitempty"`
}

func (l *ListMeta) HasMore() bool {
	return l.After != ""
}

// GetListMeta returns a ListMeta struct (itself). It exists because any
// structs that embed ListMeta will inherit it, and thus implement the
// ListContainer interface.
func (l *ListMeta) GetListMeta() *ListMeta {
	return l
}
