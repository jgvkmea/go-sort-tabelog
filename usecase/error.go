package usecase

import "errors"

var (
	BadRequestError = errors.New("bad request error")
	InternalError   = errors.New("internal request error")
)
