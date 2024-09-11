package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-chi-boilerplate/src/config"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

type (
	GoMiddleware interface {
		LogRequest(next http.Handler) http.Handler
		RecoverPanic(next http.Handler) http.Handler
		BasicAuth(username, password string) func(http.Handler) http.Handler
	}

	GoMiddlewareImpl struct {
		Config config.Config
	}
)

const (
	ParamQueryPage    = "page"
	ParamQueryLimit   = "limit"
	ParamQueryOffset  = "offset"
	ParamQueryKeyword = "keyword"
)

func InitMiddleware(cfg config.Config) GoMiddleware {
	return &GoMiddlewareImpl{
		Config: cfg,
	}
}

func (m *GoMiddlewareImpl) LogRequest(next http.Handler) http.Handler {
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
	reqId := "-"
	if requestId := rHeader.Get("request-id"); requestId != "" {
		reqId = requestId
	}

	return fmt.Sprintf("[IN_REQUEST: [%s] %s] REQUEST_ID: %s HEADER: %s", r.Method, r.URL.String(), reqId, string(headerByte))
}

func ServerError(w http.ResponseWriter, err error, code int) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Output(2, trace)
	w.Header().Set("Content-Type", "application/json")

	http.Error(w, http.StatusText(code), code)
}

func (m *GoMiddlewareImpl) RecoverPanic(next http.Handler) http.Handler {
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

func (m *GoMiddlewareImpl) BasicAuth(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			}

			// Decode the Base64 encoded auth string
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || parts[0] != "Basic" {
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			}

			userPass := strings.SplitN(string(decoded), ":", 2)
			if len(userPass) != 2 || userPass[0] != username || userPass[1] != password {
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
