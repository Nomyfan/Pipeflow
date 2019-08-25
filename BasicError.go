package pipeflow

// BasicError is used to raise error
type BasicError struct {
	Message string
}

func (e BasicError) Error() string {
	return e.Message
}
