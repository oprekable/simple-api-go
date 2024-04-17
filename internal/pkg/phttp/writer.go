package phttp

import (
	"fmt"
	"net/http"
	"reflect"
	"simple-api-go/internal/pkg/phttp/variable"

	"github.com/segmentio/encoding/json"
)

// HTTPHandlerContext ...
type HTTPHandlerContext struct {
	M interface{}
	E map[error]variable.ErrorResponse
}

// NewContextHandler ...
func NewContextHandler(meta interface{}) HTTPHandlerContext {
	var errMap = map[error]variable.ErrorResponse{}

	return HTTPHandlerContext{
		M: meta,
		E: errMap,
	}
}

// AddError ...
func (hhc HTTPHandlerContext) AddError(key error, value variable.ErrorResponse) {
	hhc.E[key] = value
}

// AddAndMapErrors ...
func (hhc HTTPHandlerContext) AddAndMapErrors(errs []*error, responses []variable.ResponseCode) {
	for i := range errs {
		errResponse := appendErrorsMapping(*errs[i], responses)
		hhc.AddError(*errs[i], errResponse)
	}
}

// AddErrorMap ...
func (hhc HTTPHandlerContext) AddErrorMap(errMap map[error]variable.ErrorResponse) {
	for k, v := range errMap {
		hhc.E[k] = v
	}
}

func appendErrorsMapping(err error, responses []variable.ResponseCode) variable.ErrorResponse {
	errorCode := err.Error()
	response := variable.ErrUnknown.Response

	httpStatus := http.StatusInternalServerError

	if len(responses) > 0 {
	bh:
		for i := range responses {
			if responses[i].ErrorCode == errorCode {
				errorMessages := responses[i].Messages
				if len(errorMessages) == 1 {
					response = variable.Response{
						ResponseCode: errorCode,
						ResponseDesc: &variable.ResponseDesc{
							ID: responses[i].Messages[0].ID,
							EN: responses[i].Messages[0].EN,
						},
					}

					httpStatus = responses[i].StatusCode
					break bh
				}
			}
		}
	}

	return variable.ErrorResponse{
		Response:   response,
		HTTPStatus: httpStatus,
	}
}

// CustomWriter ...
type CustomWriter struct {
	C HTTPHandlerContext
}

// Write ...
func (c *CustomWriter) Write(w http.ResponseWriter, data interface{}, nextPage string, httpStatus int, responseCode string) {
	wContentType := w.Header().Get("Content-Type")
	if wContentType != "" {
		writeResponseRaw(w, fmt.Sprintf("%v", data), wContentType, httpStatus)
		return
	}

	var response variable.SuccessResponse
	voData := reflect.ValueOf(data)
	arrayData := make([]interface{}, 0)

	if voData.Kind() != reflect.Slice {
		dataString := fmt.Sprintf("%v", data)
		if voData.Kind() == reflect.String && isJSON(dataString) {
			writeResponseRaw(w, dataString, "application/json", httpStatus)
			return
		}

		if voData.IsValid() {
			arrayData = []interface{}{data}
		}
		response.Data = arrayData
	} else {
		if voData.Len() != 0 {
			response.Data = data
		} else {
			response.Data = arrayData
		}
	}

	if httpStatus == 0 {
		httpStatus = http.StatusOK
	}

	response.ResponseCode = responseCode
	response.Next = nextPage
	response.Meta = c.C.M

	writeResponse(w, response, "application/json", httpStatus)
}

// WriteError sending error response based on err type
func (c *CustomWriter) WriteError(w http.ResponseWriter, err error) {
	var errorResponse variable.ErrorResponse
	if len(c.C.E) > 0 {
		errorResponse = LookupError(c.C.E, err)
		if errorResponse == (variable.ErrorResponse{}) {
			errorResponse = variable.ErrUnknown
		}
	} else {
		errorResponse = variable.ErrUnknown
	}

	errorResponse.Meta = c.C.M
	writeErrorResponse(w, errorResponse)
}

func writeResponse(w http.ResponseWriter, response interface{}, contentType string, httpStatus int) {
	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("failed to unmarshal"))
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(httpStatus)
	_, _ = w.Write(res)
}

func writeResponseRaw(w http.ResponseWriter, response string, contentType string, httpStatus int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(httpStatus)
	_, _ = w.Write([]byte(response))
}

func writeErrorResponse(w http.ResponseWriter, errorResponse variable.ErrorResponse) {
	writeResponse(w, errorResponse, "application/json", errorResponse.HTTPStatus)
}

// LookupError will get error message based on error type, with variables if you want give dynamic message error
func LookupError(lookup map[error]variable.ErrorResponse, err error) (res variable.ErrorResponse) {
	if msg, ok := lookup[err]; ok {
		res = msg
	}

	return
}

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
