package go_ws_sh

import "sync"

func PromiseAll(tasks []func() (interface{}, error)) *SafeChannel[any] {
	resultCh := NewSafeChannel[any]()
	go func() {
		defer (resultCh).Close()

		var wg sync.WaitGroup
		errChan := NewSafeChannel[error](len(tasks))
		results := make([]interface{}, len(tasks))

		for i, task := range tasks {
			wg.Add(1)
			go func(index int, t func() (interface{}, error)) {
				defer wg.Done()
				result, err := t()
				if err != nil {
					errChan.Send(err)
					return
				}
				results[index] = result
			}(i, task)
		}

		go func() {
			wg.Wait()
			(errChan).Close()
		}()

		var firstErr error
		for {

			var err, ok = errChan.Receive()
			if !ok {
				break
			}
			if err != nil && firstErr == nil {
				firstErr = err
			}
		}

		if firstErr != nil {

			resultCh.Send(firstErr)
		} else {
			resultCh.Send(results)

		}
	}()
	return resultCh
}
