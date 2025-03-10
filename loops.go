package logger

import (
	"time"
)

// InfiniteLoop continuously executes a provided function until it produces a matching error or a specified exception.
func (l *Logger) InfiniteLoop(f func() error, exception ...string) {
	for {
		err := f()
		if err != nil {
			if Is(err, exception...) {
				_ = l.ErrorF("infinite loop stopped: %s", err.Error())
				return
			}
		}

		time.Sleep(l.cnf.Wait)
	}
}

// Retry executes a function repeatedly until it succeeds or the maximum retries are reached,
// or a matching error is encountered.
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

// InfiniteRetry executes a function repeatedly until it succeeds or returns an error containing specified substrings.
// If a matching error occurs, it is returned immediately. Otherwise, it waits for a configured duration before retrying.
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

// InfiniteRetryIfXError retries a function f() infinitely only if the error returned by f() matches "xError".
// If f() returns nil (no error), the function exits successfully.
// If f() returns an error that is different from "xError", it stops immediately and returns that error.
//
// Additionally, to avoid log saturation:
// - It prints the "xError" message only ONCE, the first time it occurs.
// - If the issue is later resolved (f() stops returning an error), it logs that the problem was fixed.
func (l *Logger) InfiniteRetryIfXError(f func() error, exception string) error {
	var xErrorWasLogged bool // track whether we've logged xError yet

	for {
		err := f()
		if err != nil && Is(err, exception) {
			if !xErrorWasLogged {
				_ = l.ErrorF("An error occurred (%s), will keep retrying indefinitely...", err.Error())
				xErrorWasLogged = true
			}
			time.Sleep(l.cnf.Wait)
			continue
		}

		return err
	}
}
