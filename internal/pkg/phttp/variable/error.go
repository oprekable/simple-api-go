package variable

import "net/http"

// ErrUnknown ...
var ErrUnknown = ErrorResponse{
	Response: Response{
		ResponseCode: "100000",
		ResponseDesc: &ResponseDesc{
			ID: "error tidak diketahui, silahkan coba lagi",
			EN: "unknown error, please retry",
		},
	},
	HTTPStatus: http.StatusInternalServerError,
}
