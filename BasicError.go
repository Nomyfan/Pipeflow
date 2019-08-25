package pipeflow

type BasicError struct {
	Message string
}

func (e BasicError) Error() string {
	return e.Message
}
