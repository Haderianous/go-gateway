package gateway

type Validatable interface {
	Validate(localize Language) (data any, err error, errors map[string]any)
}
