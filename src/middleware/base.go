package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-chi-boilerplate/src/config"
	"io"
	"log"
	"net/http"
	"runtime/debug"
)

type (
	GoMiddleware struct {
		Config config.Config
	}
)

const (
	ParamQueryPage    = "page"
	ParamQueryLimit   = "limit"
	ParamQueryOffset  = "offset"
	ParamQueryKeyword = "keyword"
)

func InitMiddleware(cfg config.Config) *GoMiddleware {
	return &GoMiddleware{
		Config: cfg,
	}
}

func (m *GoMiddleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestLog := MapLogRequest(r)
		fmt.Println(requestLog)

		next.ServeHTTP(w, r)
	})
}

// MapLogRequest for map log request
func MapLogRequest(r *http.Request) string {
	rHeader := r.Header
	headerByte, _ := json.Marshal(rHeader)

	if rHeader.Get("Content-Type") == "application/json" {
		// Read the content
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
		}
		// Use the content
		var req interface{}
		json.Unmarshal(bodyBytes, &req)
		bodyBytes, _ = json.Marshal(req)

		// Restore the io.ReadCloser to its original state
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	return fmt.Sprint("[IN_REQUEST: ", r.URL, "] REQUEST_ID: ", rHeader.Get("request-id"), " HEADER:", string(headerByte))
}

func ServerError(w http.ResponseWriter, err error, code int) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Output(2, trace)
	w.Header().Set("Content-Type", "application/json")

	http.Error(w, http.StatusText(code), code)
}

func (m *GoMiddleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				ServerError(w, fmt.Errorf("%s", err), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
