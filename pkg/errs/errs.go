package errs

import "errors"

var (
	ErrValidation=errors.New("validation error")
	ErrUserNotFound=errors.New("user not found")
	
)