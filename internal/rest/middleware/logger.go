package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/HunterGooD/go_test_task/internal/domain/interfaces"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	log interfaces.Logger
}

// bodyWriter for response body
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func NewMiddleware(logger interfaces.Logger) *Middleware {
	return &Middleware{logger}
}

func (m *Middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		contentType := c.GetHeader("Content-Type")
		if strings.HasPrefix(contentType, "application/json") {
			// read request body
			reqBody, err := io.ReadAll(c.Request.Body)
			if err != nil {
				m.log.Warn("Failed to read request body", map[string]any{
					"error": err.Error(),
				})
			}
			// repeat read in handler
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			m.log.Info("Request received", map[string]any{
				"method":       c.Request.Method,
				"path":         c.FullPath(),
				"req_path":     c.Request.URL.Path,
				"client":       c.ClientIP(),
				"query_params": c.Request.URL.RawQuery,
			})
			m.log.Debug("Request received", map[string]any{
				"method":       c.Request.Method,
				"path":         c.FullPath(),
				"client":       c.ClientIP(),
				"req_path":     c.Request.URL.Path,
				"query_params": c.Request.URL.RawQuery,
				"req_body":     string(reqBody),
			})
		} else {
			m.log.Info("Request received (non-JSON body)", map[string]any{
				"method":       c.Request.Method,
				"path":         c.FullPath(),
				"req_path":     c.Request.URL.Path,
				"client":       c.ClientIP(),
				"query_params": c.Request.URL.RawQuery,
				"content_type": contentType,
			})
		}

		// write response data to writer
		writer := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		if c.Writer.Status() == http.StatusNotFound {
			m.log.Warn("Route not found", map[string]any{
				"method":       c.Request.Method,
				"path":         c.FullPath(),
				"req_path":     c.Request.URL.Path,
				"client":       c.ClientIP(),
				"query_params": c.Request.URL.RawQuery,
			})
		}

		resContentType := c.Writer.Header().Get("Content-Type")
		latency := time.Since(start)
		if strings.HasPrefix(resContentType, "application/json") {
			m.log.Info("Response sent", map[string]any{
				"status":      http.StatusText(c.Writer.Status()),
				"status_code": c.Writer.Status(),
				"client":      c.ClientIP(),
				"latency":     latency,
			})
			m.log.Debug("Response sent", map[string]any{
				"status":      http.StatusText(c.Writer.Status()),
				"status_code": c.Writer.Status(),
				"client":      c.ClientIP(),
				"res_body":    writer.body.String(),
				"latency":     latency,
			})
		} else {
			m.log.Info("Response sent (non-JSON body)", map[string]any{
				"status":       http.StatusText(c.Writer.Status()),
				"status_code":  c.Writer.Status(),
				"client":       c.ClientIP(),
				"content_type": resContentType,
				"latency":      latency,
			})
		}

		// TODO: search methods for example
		// slogStrings := make([]slog.Attr, 0)
		// slogStrings = append(slogStrings, slog.String("asd", "asdas"))
		// slogStrings = append(slogStrings, slog.String("asd", "asdas"))
		// slog.Info("TODO: use this method", slogStrings...)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				m.log.Error("Request error", map[string]any{
					"error": e.Error(),
				})
			}
		}

	}
}
