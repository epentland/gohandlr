package options

import "net/http"

type Writer interface {
	Write(w http.ResponseWriter, r *http.Request, buff any) error
	Accept() string
}
