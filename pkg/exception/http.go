package exception

import "fmt"

type HttpError struct {
	Code    int               `json:"-"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func (err HttpError) Error() string {
	return err.Message
}

func Catch(err error) {
	if err != nil {
		panic(err)
	}
}

func InternalError(message string) error {
	return HttpError{
		Code:    500,
		Message: message,
	}
}

func Http(code int, message string) error {
	return HttpError{
		Code:    code,
		Message: message,
	}
}

func BadRequest(message string) error {
	return HttpError{
		Code:    400,
		Message: message,
	}
}

func Forbidden() error {
	return HttpError{
		Code:    403,
		Message: "Forbidden",
	}
}

func NotFound(field string) error {
	return HttpError{
		Code:    404,
		Message: fmt.Sprintf("%s tidak ditemukan", field),
	}
}
