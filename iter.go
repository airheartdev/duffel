package duffel

type Iter[T any] struct {
	cur      *T
	err      error
	list     ListContainer[T]
	meta     *ListMeta
	nextPage PageFn[T]
	values   []*T
}

func Collect[T any](it *Iter[T]) ([]*T, error) {
	collection := make([]*T, 0)
	for it.Next() {
		collection = append(collection, it.Current())
	}
	return collection, it.Err()
}

type PageFn[T any] func(meta *ListMeta) (*List[T], error)

// Current returns the most recent item
// visited by a call to Next.
func (it *Iter[T]) Current() *T {
	if it == nil {
		return nil
	}

	return it.cur
}

// Err returns the error, if any,
// that caused the Iter to stop.
// It must be inspected
// after Next returns false.
func (it *Iter[T]) Err() error {
	if it == nil {
		return nil
	}
	return it.err
}

// List returns the current list object which the iterator is currently using.
// List objects will change as new API calls are made to continue pagination.
func (it *Iter[T]) List() ListContainer[T] {
	return it.list
}

// Meta returns the list metadata.
func (it *Iter[T]) Meta() *ListMeta {
	if it == nil {
		return nil
	}
	return it.meta
}

// Next advances the Iter to the next item in the list,
// which will then be available
// through the Current method.
// It returns false when the iterator stops
// at the end of the list.
func (it *Iter[T]) Next() bool {
	if it.err != nil {
		return false
	}

	if len(it.values) == 0 && it.meta.HasMore() {
		it.getPage()
	}

	if len(it.values) == 0 {
		return false
	}
	it.cur = it.values[0]
	it.values = it.values[1:]
	return true
}

func (it *Iter[T]) getPage() {
	it.list, it.err = it.nextPage(it.meta)
	if it.err == nil {
		it.values = it.list.GetItems()
		it.meta = it.list.GetListMeta()
	}
}

// GetIter returns a new Iter for a given query and type.
func GetIter[T any](pager PageFn[T]) *Iter[T] {
	iter := &Iter[T]{
		nextPage: pager,
		meta:     &ListMeta{},
	}

	iter.getPage()

	return iter
}
