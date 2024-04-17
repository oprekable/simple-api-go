package log

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"runtime"
	"strings"
	"time"

	"github.com/DmitriyVTitov/size"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

const (
	// CRequest ...
	CRequest string = "request"

	// CStatus ...
	CStatus string = "status"

	// CSize ...
	CSize string = "size"

	// CDuration ...
	CDuration string = "duration"

	// CResponseBody ...
	CResponseBody string = "response_body"

	// CResponseHeader ...
	CResponseHeader string = "response_header-"

	// CHttpURL ...
	CHttpURL string = "http_url"

	// CHttpMethod ...
	CHttpMethod string = "http_method"

	// CIp ...
	CIp string = "ip"

	// CUserAgent ...
	CUserAgent string = "user_agent"

	// CRealIP ...
	CRealIP string = "real_ip"

	// CReferer ...
	CReferer string = "referer"

	// CUpTime ...
	CUpTime string = "uptime"
)

type LogCaller struct {
	File     string
	Function string
	Line     int
}

func NewLogCaller() LogCaller {
	pc, f, l, _ := runtime.Caller(4)
	return LogCaller{
		File:     f,
		Line:     l,
		Function: runtime.FuncForPC(pc).Name(),
	}
}

func (lc LogCaller) MarshalZerologObject(e *zerolog.Event) {
	e.Str("file", lc.File).
		Int("line", lc.Line).
		Str("function", lc.Function)
}

// NewHTTPRequestLogger ...
func NewHTTPRequestLogger(mainCtx context.Context, logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return alice.New(hlog.NewHandler(logger),
			hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
				hlog.FromRequest(r).
					Info().
					Ctx(mainCtx).
					Int(CStatus, status).
					Int(CSize, size).
					Str(CDuration, duration.String()).
					Str(middleware.RequestIDHeader, middleware.GetReqID(r.Context())).
					Msg("")
			}),
			hlog.URLHandler(CHttpURL),
			hlog.MethodHandler(CHttpMethod),
			hlog.RemoteAddrHandler(CIp),
			hlog.UserAgentHandler(CUserAgent),
			hlog.CustomHeaderHandler(CRealIP, "X-Real-Ip"),
			hlog.RefererHandler(CReferer)).
			Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				logRedefine := zerolog.Ctx(r.Context())
				requestDump, err := httputil.DumpRequest(r, true)

				if err != nil {
					logRedefine.UpdateContext(func(c zerolog.Context) zerolog.Context {
						return c.Object("caller", NewLogCaller()).Err(err)
					})
				}

				if size.Of(requestDump) < 10000 && !strings.Contains(string(requestDump), "password") {
					logRedefine.UpdateContext(func(c zerolog.Context) zerolog.Context {
						return c.Str(CRequest, string(requestDump))
					})
				}

				responseRecorder := httptest.NewRecorder()
				next.ServeHTTP(responseRecorder, r)

				for k, v := range responseRecorder.Header() {
					w.Header().Set(k, strings.Join(v, ""))
					logRedefine.UpdateContext(func(c zerolog.Context) zerolog.Context {
						return c.Strs(CResponseHeader+k, v)
					})
				}

				w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(r.Context()))

				w.WriteHeader(responseRecorder.Code)
				responseBody := responseRecorder.Body.String()

				// GZIP decode
				if len(responseRecorder.Header().Get("Content-Encoding")) > 0 && responseRecorder.Header().Get("Content-Encoding") == "gzip" {
					reader := bytes.NewReader(responseRecorder.Body.Bytes())
					gzipReader, err := gzip.NewReader(reader)

					if err != nil {
						logRedefine.UpdateContext(func(c zerolog.Context) zerolog.Context {
							return c.Object("caller", NewLogCaller()).Err(err)
						})
					}

					unzipOutput, err := io.ReadAll(gzipReader)

					if err != nil {
						logRedefine.UpdateContext(func(c zerolog.Context) zerolog.Context {
							return c.Object("caller", NewLogCaller()).Err(err)
						})
					}

					if len(unzipOutput) > 0 {
						responseBody = string(unzipOutput)
					}
				}

				_, err = responseRecorder.Body.WriteTo(w)
				if err != nil {
					logRedefine.UpdateContext(func(c zerolog.Context) zerolog.Context {
						return c.Object("caller", NewLogCaller()).Err(err)
					})
				}

				if size.Of(responseBody) < 10000 {
					logRedefine.UpdateContext(func(c zerolog.Context) zerolog.Context {
						return c.Str(CResponseBody, responseBody)
					})
				}
			}))
	}
}
