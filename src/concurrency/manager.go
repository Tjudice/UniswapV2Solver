package concurrency

import "golang.org/x/sync/errgroup"

type ConcurrentCaller[T any] struct {
	ret   chan T
	limit int
	calls []func() (T, error)
}

func NewConcurrentCaller[T any](limit int) *ConcurrentCaller[T] {
	if limit == 0 {
		limit = 10
	}
	return &ConcurrentCaller[T]{
		ret:   make(chan T, 1),
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
	c.ret = make(chan T, len(c.calls))
	for _, call := range c.calls {
		cc := call
		wg.Go(func() error {
			ret, err := cc()
			if err != nil {
				return err
			}
			c.ret <- ret
			return nil
		})
	}
	err := wg.Wait()
	if err != nil {
		return nil, err
	}
	r := []T{}
	for x := range c.ret {
		r = append(r, x)
	}
	return r, nil
}
