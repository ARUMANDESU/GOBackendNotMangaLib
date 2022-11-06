package models

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

var ErrNoRecord = errors.New("model: no matching record found")
var TokenError = errors.New("server Error")
var ExpiredToken = errors.New("token is expired")
var NotAuthorized = errors.New("you are not authorized")

type ErrorHandlerJwt struct {
	Payload jwt.MapClaims
	Err     error
}

func HandleJWTError(payload jwt.MapClaims, err error) *ErrorHandlerJwt {
	return &ErrorHandlerJwt{
		Err:     err,
		Payload: payload,
	}
}
