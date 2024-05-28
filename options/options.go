package options

type Options interface{}

func DefaultOptions() []Options {
	return []Options{
		WithJSONBodyReader(),
		WithParamsReader(),
	}
}
