package core

import "errors"

// General Error
var (
	ErrInternal       = errors.New("100001")
	ErrDBConn         = errors.New("100002")
	ErrUnauthorized   = errors.New("100003")
	ErrInvalidHeader  = errors.New("100004")
	ErrInvalidPayload = errors.New("100005")
	ErrRedisConn      = errors.New("100006")
	ErrDataNotFound   = errors.New("100007")

	// ErrEndpointDeprecated error type for Error of Endpoint Not Available (Deprecated)
	ErrEndpointDeprecated = errors.New("100008")
	ErrBadPayload         = errors.New("100009")

	ErrSignatureExpired  = errors.New("100010")
	ErrFileSizeOverLimit = errors.New("100011")
)
