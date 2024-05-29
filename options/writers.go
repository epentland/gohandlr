package options

import "net/http"

type Writer interface {
	Write(http.ResponseWriter, *http.Request, any) error
	Accept() string
}
