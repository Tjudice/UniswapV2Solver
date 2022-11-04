package concurrency

import (
	"log"
	"sync"

	"golang.org/x/sync/errgroup"
)

type ConcurrentCaller[T any] struct {
	mut   sync.Mutex
	limit int
	calls []func() (T, error)
}

func NewConcurrentCaller[T any](limit int) *ConcurrentCaller[T] {
	if limit == 0 {
		limit = 10
	}
	return &ConcurrentCaller[T]{
		mut:   sync.Mutex{},
		limit: limit,
		calls: make([]func() (T, error), 0),
	}
}

func (c *ConcurrentCaller[T]) AddCall(call func() (T, error)) {
	c.calls = append(c.calls, call)
}

func (c *ConcurrentCaller[T]) Run() ([]T, error) {
	wg := errgroup.Group{}
	wg.SetLimit(c.limit)
	r := []T{}
	for _, call := range c.calls {
		log.Println("calling")
		cc := call
		wg.Go(func() error {
			ret, err := cc()
			if err != nil {
				return err
			}
			c.mut.Lock()
			r = append(r, ret)
			c.mut.Unlock()
			return nil
		})
	}
	err := wg.Wait()
	if err != nil {
		return nil, err
	}
	return r, nil
}
