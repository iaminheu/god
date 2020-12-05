package httpx

const defaultCode = 1001

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func (e *CodeError) Error() string {
	return e.Msg
}

func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

func NewDefaultError(msg string) error {
	return NewCodeError(defaultCode, msg)
}
