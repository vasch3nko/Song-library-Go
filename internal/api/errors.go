package api

import "strconv"

type HttpError struct {
    StatusCode int
}

func (e *HttpError) Error() string {
    return strconv.Itoa(e.StatusCode)
}

func NewHttpError(statusCode int) *HttpError {
    return &HttpError{StatusCode: statusCode}
}
