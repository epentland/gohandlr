package options

type Options interface{}

func WithDefaults() []Options {
	return []Options{
		WithJSONBodyReader(),
		WithParamsReader(),
	}
}
