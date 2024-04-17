package error

import (
	"simple-api-go/internal/app/error/core"
)

// Errors Register new errors here!
var Errors = []*error{
	&core.ErrInternal,
	&core.ErrDBConn,
	&core.ErrUnauthorized,
	&core.ErrInvalidHeader,
	&core.ErrInvalidPayload,
	&core.ErrRedisConn,
	&core.ErrDataNotFound,
	&core.ErrEndpointDeprecated,
	&core.ErrBadPayload,
	&core.ErrSignatureExpired,
	&core.ErrFileSizeOverLimit,
}
