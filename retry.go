package logger

func (l *Logger) RunWithRetries(status int, f func() error) *Error {
	var retries = 0
	for {
		err := f()
		if err != nil {
			retries++
			e := l.ErrorF(status, err.Error())
			if retries >= l.cnf.Retries {
				return e
			}
		} else {
			return nil
		}
	}
}
