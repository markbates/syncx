package syncx

import (
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

func (g *Group) Wait() {
	if g == nil {
		return
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
}

func (g *Group) report(err error) {
	if g == nil {
		return
	}

	if err == nil {
		return
	}

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
