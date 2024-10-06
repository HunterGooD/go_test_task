package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// bodyWriter for response body
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		contentType := c.GetHeader("Content-Type")
		if strings.HasPrefix(contentType, "application/json") {
			// read request body
			reqBody, err := io.ReadAll(c.Request.Body)
			if err != nil {
				slog.Warn("Failed to read request body", slog.String("error", err.Error()))
			}
			// repeat read in handler
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			slog.Info("Request received",
				slog.String("method", c.Request.Method),
				slog.String("path", c.FullPath()),
				slog.String("req_path", c.Request.URL.Path),
				slog.String("client", c.ClientIP()),
				slog.String("query_params", c.Request.URL.RawQuery),
			)
			slog.Debug("Request received",
				slog.String("client", c.ClientIP()),
				slog.String("req_path", c.Request.URL.Path),
				slog.String("req_body", string(reqBody)),
			)
		} else {
			slog.Info("Request received (non-JSON body)",
				slog.String("method", c.Request.Method),
				slog.String("path", c.FullPath()),
				slog.String("req_path", c.Request.URL.Path),
				slog.String("client", c.ClientIP()),
				slog.String("query_params", c.Request.URL.RawQuery),
				slog.String("content_type", contentType),
			)
		}

		// write response data to writer
		writer := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		if c.Writer.Status() == http.StatusNotFound {
			slog.Warn("Route not found",
				slog.String("method", c.Request.Method),
				slog.String("path", c.FullPath()),
				slog.String("req_path", c.Request.URL.Path),
				slog.String("client", c.ClientIP()),
				slog.String("query_params", c.Request.URL.RawQuery),
			)
		}

		resContentType := c.Writer.Header().Get("Content-Type")
		latency := time.Since(start)
		if strings.HasPrefix(resContentType, "application/json") {
			slog.Info("Response sent",
				slog.String("status", http.StatusText(c.Writer.Status())),
				slog.Int("status_code", c.Writer.Status()),
				slog.String("client", c.ClientIP()),
				slog.Duration("latency", latency),
			)
			slog.Debug("Response sent",
				slog.String("status", http.StatusText(c.Writer.Status())),
				slog.Int("status_code", c.Writer.Status()),
				slog.String("client", c.ClientIP()),
				slog.String("res_body", writer.body.String()),
				slog.Duration("latency", latency),
			)
		} else {
			slog.Info("Response sent (non-JSON body)",
				slog.String("status", http.StatusText(c.Writer.Status())),
				slog.Int("status_code", c.Writer.Status()),
				slog.String("client", c.ClientIP()),
				slog.String("content_type", resContentType),
				slog.Duration("latency", latency),
			)
		}

		// TODO: search methods for example
		// slogStrings := make([]slog.Attr, 0)
		// slogStrings = append(slogStrings, slog.String("asd", "asdas"))
		// slogStrings = append(slogStrings, slog.String("asd", "asdas"))
		// slog.Info("TODO: use this method", slogStrings...)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				slog.Error("Request error", slog.String("error", e.Error()))
			}
		}

	}
}
