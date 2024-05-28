package handle

import "net/http"

type Error interface {
	Error() string
	Status() int
}

type NewError struct {
	err    error
	status int
}

func (e NewError) Error() string {
	return e.err.Error()
}

func (e NewError) Status() int {
	return e.status
}

func ErrorInternal(err error) Error {
	return NewError{err: err, status: http.StatusInternalServerError}
}

func ErrorBadGateway(err error) Error {
	return NewError{err: err, status: http.StatusBadGateway}
}

func ErrorUnavailable(err error) Error {
	return NewError{err: err, status: http.StatusServiceUnavailable}
}

func ErrorTimeout(err error) Error {
	return NewError{err: err, status: http.StatusGatewayTimeout}
}
