package logger

import (
	"time"
)

// InfiniteLoop executes a function infinitely until it returns an 
// error that contains any of the specified substrings. 
// It waits for the configured time between each iteration.
func (l *Logger) InfiniteLoop(f func() error, exception ...string) {
	for {
		err := f()
		if err != nil {
			if Is(err, exception...) {
				l.ErrorF("infinite loop stopped: %s", err.Error())
				return
			}
		}

		time.Sleep(l.cnf.Wait)
	}
}

// Retry executes a function and retries if any error. 
// It stops after the number of retries specified in the configuration or if 
// the error contains any of the specified substrings.
// It waits for the configured time between each iteration.
func (l *Logger) Retry(f func() error, exception ...string) error {
	var retries = 0
	for {
		err := f()
		if err != nil {
			retries++

			if Is(err, exception...) {
				return err
			}

			if retries >= l.cnf.Retries {
				return err
			}

			time.Sleep(l.cnf.Wait)
		} else {
			return nil
		}
	}
}

// InfiniteRetry executes a function and retries infinitely if any error. 
// It stops if the error contains any of the specified substrings.
// It waits for the configured time between each iteration.
func (l *Logger) InfiniteRetry(f func() error, exception ...string) error {
	for {
		err := f()
		if err != nil {
			if Is(err, exception...) {
				return err
			}

			time.Sleep(l.cnf.Wait)
		} else {
			return nil
		}
	}
}
