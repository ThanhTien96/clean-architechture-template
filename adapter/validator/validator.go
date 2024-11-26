package validator

type Validator interface {
	Validate(interface{}) error
	Message() []string
}
