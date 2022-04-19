package syncx

import (
	"fmt"
	"sync"
)

// Group is similar to errgroup.Group, but it
// gives you an errors channel to listen to
// all of the errors from the goroutines.
type Group struct {
	mu        sync.RWMutex
	wg        sync.WaitGroup
	eg        sync.WaitGroup
	stop      sync.Once
	err       error // first error
	listeners map[int]chan error
}

func (g *Group) Go(fn func() error) {
	if g == nil {
		return
	}

	if fn == nil {
		return
	}

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		err := fn()
		if err != nil {
			g.report(err)
		}

	}()
}

// Errors returns a channel that will receive all of the errors
// from the goroutines. This channel will be closed when the group
// is done.
// This method should be called **BEFORE** loading any Go functions.
// otherwise you run the risk of missing errors.
// 		wg := &Group{}
// 		errs := wg.Errors()
// 		wg.Go(func() error {
// 			return nil
// 		})
// 		wg.Wait()
func (g *Group) Errors() <-chan error {
	if g == nil {
		return nil
	}

	ch := make(chan error)

	g.mu.Lock()
	if g.listeners == nil {
		g.listeners = map[int]chan error{}
	}
	g.listeners[len(g.listeners)] = ch
	g.mu.Unlock()

	return ch
}

// Wait blocks until all of the goroutines have completed.
// If any of the goroutines return an error, Wait will return
// the first error reported.
// To get all of the errors, use `Errors()`
func (g *Group) Wait() error {
	if g == nil {
		return fmt.Errorf("group is nil")
	}

	defer func() {
		g.stop.Do(func() {

			g.eg.Wait()

			for i, ch := range g.listeners {
				g.mu.Lock()
				close(ch)
				delete(g.listeners, i)
				g.mu.Unlock()
			}

		})

	}()

	g.wg.Wait()

	return g.err
}

func (g *Group) report(err error) {
	if g == nil {
		return
	}

	if err == nil {
		return
	}

	g.mu.Lock()
	if g.err == nil {
		g.err = err
	}
	g.mu.Unlock()

	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, ch := range g.listeners {
		g.eg.Add(1)
		go func(ch chan error) {
			defer g.eg.Done()
			ch <- err
		}(ch)
	}
}
