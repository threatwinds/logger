package logger

import (
	"strings"
	"time"
)

// RunWithRetries executes a function and retries if any error.
func (l *Logger) RunWithRetries(f func() error, exception ...string) *Error {
	var retries = 0
	for {
		err := f()
		if err != nil {
			retries++
			e := l.ErrorF(err.Error())

			for _, ex := range exception {
				if strings.Contains(err.Error(), ex) {
					return e
				}
			}
			
			if retries >= l.cnf.Retries {
				return e
			}

			time.Sleep(l.cnf.Wait)
		} else {
			return nil
		}
	}
}

// RunWithInfRetries executes a function and retries infinitely if any error.
func (l *Logger) RunWithInfRetries(f func() error, exception ...string) *Error {
	var retries = 0
	for {
		err := f()
		if err != nil {
			retries++
			
			e := l.ErrorF(err.Error())

			for _, ex := range exception {
				if strings.Contains(err.Error(), ex) {
					return e
				}
			}

			time.Sleep(l.cnf.Wait)
		} else {
			return nil
		}
	}
}