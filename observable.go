package main

type Observable[T any] struct {
	data        T
	subscribers []Subscriber[T]
}

type Subscriber[T any] interface {
	OnChange(old *T, val *T)
}

func NewObservable[T any](initialValue T) *Observable[T] {
	return &Observable[T]{
		data:        initialValue,
		subscribers: make([]Subscriber[T], 0),
	}
}

func (o *Observable[T]) Subscribe(s Subscriber[T]) {
	o.subscribers = append(o.subscribers, s)
}

func (o *Observable[T]) Close() {
	o.subscribers = make([]Subscriber[T], 0)
}

func (o *Observable[T]) Set(value T) {
	oldValue := o.data
	o.data = value
	for _, s := range o.subscribers {
		s.OnChange(&oldValue, &value)
	}
}

func (o *Observable[T]) Update(updater func(oldValue *T) T) {
	oldValue := o.data
	o.data = updater(&oldValue)
	for _, s := range o.subscribers {
		s.OnChange(&oldValue, &o.data)
	}
}
