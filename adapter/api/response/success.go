package response

type Success struct {
	statusCode int
	result interface{}
}

func NewSuccess(result interface{}, status int) *Success {
	return &Success{
		statusCode: status,
		result: result,
	}
}