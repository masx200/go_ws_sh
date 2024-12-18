package go_ws_sh

import "sync"

func PromiseAll(tasks []func() (interface{}, error)) <-chan interface{} {
	resultCh := make(chan interface{})
	go func() {
		defer close(resultCh)

		var wg sync.WaitGroup
		errChan := make(chan error, len(tasks))
		results := make([]interface{}, len(tasks))

		for i, task := range tasks {
			wg.Add(1)
			go func(index int, t func() (interface{}, error)) {
				defer wg.Done()
				result, err := t()
				if err != nil {
					select {
					case errChan <- err:
					default:
					}
					return
				}
				results[index] = result
			}(i, task)
		}

		go func() {
			wg.Wait()
			close(errChan)
		}()

		var firstErr error
		for err := range errChan {
			if err != nil && firstErr == nil {
				firstErr = err
			}
		}

		if firstErr != nil {
			resultCh <- firstErr
		} else {
			resultCh <- results
		}
	}()
	return resultCh
}
