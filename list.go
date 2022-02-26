package duffel

// ListContainer is a general interface for which all list object structs
// should comply. They achieve this by embedding a ListMeta struct and
// inheriting its implementation of this interface.
type ListContainer[T any] interface {
	SetListMeta(*ListMeta)
	GetListMeta() *ListMeta
	GetItems() []*T
	SetItems(items []*T)
}

type List[T any] struct {
	items []*T
	// *Iter[T]
	*ListMeta
}

func (l *List[T]) GetItems() []*T {
	return l.items
}

func (l *List[T]) SetItems(items []*T) {
	l.items = items
}

func (l *List[T]) SetListMeta(meta *ListMeta) {
	l.ListMeta = meta
}

// ListMeta is the structure that contains the common properties
// of List iterators.
type ListMeta struct {
	// HasMore is a boolean that indicates whether there are more items

	// After is a string that contains the token for the next page of results
	After string `json:"after" url:"after,omitempty"`

	// Before is a string that contains the token for the previous page of results
	Before string `json:"before" url:"-"`

	// Limit is a number that indicates the maximum number of items to return
	Limit int `json:"limit" url:"limit,omitempty"`
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
