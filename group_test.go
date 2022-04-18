package syncx

import (
	"fmt"
	"sync"
	"testing"
)

// func Test_Group(t *testing.T) {
// 	t.Parallel()

// 	N := 10

// 	wg := &Group{}

// 	var mu sync.Mutex
// 	var count int

// 	for i := 0; i < N; i++ {
// 		wg.Go(func() error {
// 			mu.Lock()
// 			count++
// 			mu.Unlock()
// 			return nil
// 		})
// 	}

// 	wg.Wait()

// 	if count != N {
// 		t.Fatalf("expected %v, got %v", N, count)
// 	}
// }

func Test_Group_Errors(t *testing.T) {
	t.Parallel()

	N := 10
	wg := &Group{}
	errs := wg.Errors()

	for i := 0; i < N; i++ {
		i := i
		wg.Go(func() error {
			return fmt.Errorf("error %v", i)
		})
	}

	quit := make(chan struct{})

	var mu sync.Mutex
	var act []error
	go func() {
		defer close(quit)

		for e := range errs {
			if e == nil {
				break
			}

			mu.Lock()
			act = append(act, e)
			mu.Unlock()

		}
		fmt.Println("done with loop")
	}()

	wg.Wait()

	<-quit

	if len(act) != N {
		t.Fatalf("expected %v, got %v", N, len(act))
	}
}
