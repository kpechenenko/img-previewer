package handler

type InvalidParamErr struct {
	description string
}

func (e *InvalidParamErr) Error() string {
	return e.description
}
