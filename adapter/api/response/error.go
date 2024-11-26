package response

type Error struct {
	statusCode int
	Errors     []string `json:"errors"`
}


func NewError(err error, status int) *Error {
	return &Error{
		statusCode: status,
		Errors: []string{err.Error()},
	}
}

func NewErrorMessage(messages []string, status int) *Error {
	return &Error{
		statusCode: status,
		Errors: messages,
	}
}
