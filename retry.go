package logger

func (l *Logger) RunWithRetries(f func() *Error) *Error {
	var retries = 0
	for {
		e := f()
		if e != nil {
			return nil
		} else {
			retries++
			if retries >= l.cnf.Retries {
				return e
			}
		}
	}
}
