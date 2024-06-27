package domain

type ExternalServiceError struct {
	Err error
}

func (e *ExternalServiceError) Error() string {
	return e.Err.Error()
}
