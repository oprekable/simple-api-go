package checkheader

import (
	"context"
	"net/http"
	"reflect"
	"simple-api-go/internal/app/error/core"
	"simple-api-go/internal/pkg/phttp"
	"simple-api-go/internal/pkg/utils/httphelper"
	"simple-api-go/internal/pkg/utils/log"
	"simple-api-go/variable"
	"strconv"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/segmentio/encoding/json"
)

const (
	CRequiredHeaders variable.KeyType = "required_headers"
	CVisitorIP       string           = "ip"
)

// CheckRequiredHeader ...
func CheckRequiredHeader(httpCtx phttp.HTTPHandlerContext, headers map[string]reflect.Kind, watermillPublisher message.Publisher, watermillPublishTopic string) func(next http.Handler) http.Handler {
	writer := phttp.CustomWriter{
		C: httpCtx,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxValue := make(map[string]interface{})

			for k, v := range headers {
				headerValue := r.Header.Get(k)
				if headerValue == "" {
					err := core.ErrInvalidHeader
					log.AddErr(r.Context(), err)
					writer.WriteError(w, err)
					return
				}

				var er error
				var convertedValue interface{}
				switch v {
				case reflect.Bool:
					{
						convertedValue, er = strconv.ParseBool(headerValue)
					}
				case reflect.Int64:
					{
						convertedValue, er = strconv.ParseInt(headerValue, 10, 64)
					}
				case reflect.Float64:
					{
						convertedValue, er = strconv.ParseFloat(headerValue, 64)
					}
				default:
					convertedValue = headerValue
				}

				if er != nil {
					err := core.ErrInvalidHeader
					log.AddErr(r.Context(), err)
					writer.WriteError(w, err)
					return
				}

				ctxValue[k] = convertedValue
			}

			ctxValue[CVisitorIP] = httphelper.GetVisitorIP(r)
			ctx := context.WithValue(r.Context(), CRequiredHeaders, ctxValue)
			r = r.WithContext(ctx)

			if watermillPublisher != nil || watermillPublishTopic != "" {
				ctxValueJsonByte, _ := json.Marshal(ctxValue)
				msg := message.NewMessage(watermill.NewUUID(), ctxValueJsonByte)
				middleware.SetCorrelationID(watermill.NewUUID(), msg)
				_ = watermillPublisher.Publish(watermillPublishTopic, msg)
			}

			next.ServeHTTP(w, r)
		})
	}
}
