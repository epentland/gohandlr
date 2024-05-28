package options

import "net/http"

type BodyReader interface {
	Reader(*http.Request, any) error
	ContentType() string
}

type ParamsReader interface {
	Reader(*http.Request, any) error
}
