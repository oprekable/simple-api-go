package log

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	var buf bytes.Buffer
	writerIO := bufio.NewWriter(&buf)
	r := chi.NewRouter()
	logger := zerolog.New(writerIO)
	r.Use(NewHTTPRequestLogger(context.Background(), logger))
	r.MethodFunc(http.MethodGet, "/test", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()
	_, _ = LogTestRequest(t, ts, "GET", "/test", nil)
	_ = writerIO.Flush() // forcefully write remaining
	assert.Regexp(t, regexp.MustCompile("^{\"level\":\"info\".*"), buf.String())
}

func LogTestRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return resp, string(respBody)
}
