package logger

import (
	"time"
)

// RunWithRetries executes a function and retries if any error.
func (l *Logger) RunWithRetries(f func() error, exception ...string) error {
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

// RunWithInfRetries executes a function and retries infinitely if any error.
func (l *Logger) RunWithInfRetries(f func() error, exception ...string) error {
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
