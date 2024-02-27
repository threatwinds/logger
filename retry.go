package logger

import "time"

// RunWithRetries executes a function and retries if any error.
func (l *Logger) RunWithRetries(status map[int][]string, f func() error) *Error {
	var retries = 0
	for {
		err := f()
		if err != nil {
			retries++
			e := l.ErrorF(status, err.Error())
			if retries >= l.cnf.Retries {
				return e
			}
			time.Sleep(l.cnf.Wait)
		} else {
			return nil
		}
	}
}
