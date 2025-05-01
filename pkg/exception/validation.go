package exception

func InvalidPayload(errors ...map[string]string) error {
	httpError := HttpError{
		Code:    422,
		Message: "Invalid payload",
	}

	if len(errors) > 0 {
		httpError.Errors = errors[0]
	}

	return httpError
}

func PasswordRequired() error {
	return InvalidPayload(map[string]string{
		"password": "Password required",
	})
}

func UserNotRegistered() error {
	return InvalidPayload(map[string]string{
		"email": "Email not registered",
	})
}

func PasswordIncorrect() error {
	return InvalidPayload(map[string]string{
		"password": "Password incorrect",
	})
}

func EmailUsed() error {
	return InvalidPayload(map[string]string{
		"email": "Email was used",
	})
}
