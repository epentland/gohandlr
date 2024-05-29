package handle

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewError(t *testing.T) {
	err := errors.New("test error")
	status := http.StatusInternalServerError
	newErr := NewError{err: err, status: status}

	if newErr.Error() != err.Error() {
		t.Errorf("Expected error message: %s, got: %s", err.Error(), newErr.Error())
	}

	if newErr.Status() != status {
		t.Errorf("Expected status code: %d, got: %d", status, newErr.Status())
	}
}

func TestErrorInternal(t *testing.T) {
	err := errors.New("internal error")
	internalErr := ErrorInternal(err)

	if internalErr.Error() != err.Error() {
		t.Errorf("Expected error message: %s, got: %s", err.Error(), internalErr.Error())
	}

	if internalErr.Status() != http.StatusInternalServerError {
		t.Errorf("Expected status code: %d, got: %d", http.StatusInternalServerError, internalErr.Status())
	}
}

func TestErrorBadGateway(t *testing.T) {
	err := errors.New("bad gateway error")
	badGatewayErr := ErrorBadGateway(err)

	if badGatewayErr.Error() != err.Error() {
		t.Errorf("Expected error message: %s, got: %s", err.Error(), badGatewayErr.Error())
	}

	if badGatewayErr.Status() != http.StatusBadGateway {
		t.Errorf("Expected status code: %d, got: %d", http.StatusBadGateway, badGatewayErr.Status())
	}
}

func TestErrorUnavailable(t *testing.T) {
	err := errors.New("service unavailable error")
	unavailableErr := ErrorUnavailable(err)

	if unavailableErr.Error() != err.Error() {
		t.Errorf("Expected error message: %s, got: %s", err.Error(), unavailableErr.Error())
	}

	if unavailableErr.Status() != http.StatusServiceUnavailable {
		t.Errorf("Expected status code: %d, got: %d", http.StatusServiceUnavailable, unavailableErr.Status())
	}
}

func TestErrorTimeout(t *testing.T) {
	err := errors.New("timeout error")
	timeoutErr := ErrorTimeout(err)

	if timeoutErr.Error() != err.Error() {
		t.Errorf("Expected error message: %s, got: %s", err.Error(), timeoutErr.Error())
	}

	if timeoutErr.Status() != http.StatusGatewayTimeout {
		t.Errorf("Expected status code: %d, got: %d", http.StatusGatewayTimeout, timeoutErr.Status())
	}
}
