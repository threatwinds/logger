package logger

import (
	"time"
)

// Retry executes a function and retries if any error.
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
