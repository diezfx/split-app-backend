package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/diezfx/split-app-backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

const httpSuccessThreshold = 399

func HTTPLoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		capturingWriter := &capturingResponseWriter{ResponseWriter: ctx.Writer}
		ctx.Writer = capturingWriter
		requestBody, _ := io.ReadAll(ctx.Request.Body)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		// Process request
		ctx.Next()

		// Log request information
		duration := time.Since(start)
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		responseBody := capturingWriter.Body()

		log := logger.Info(ctx).
			String("client_ip", clientIP).
			String("method", method)

		// only log bodies when there was an error
		if statusCode > httpSuccessThreshold {
			if json.Valid(requestBody) {
				log = log.RawJSON("request_body", requestBody)
			}
			if json.Valid(responseBody) {
				log = log.RawJSON("response_body", responseBody)
			}
		}

		log.String("client_ip", clientIP).
			String("method", method).
			String("path", path).
			Int("status_code", statusCode).
			Duration("duration", duration).
			Msg("Request")
	}
}

// Custom capturingResponseWriter to capture the response body
type capturingResponseWriter struct {
	gin.ResponseWriter
	body []byte
}

func (w *capturingResponseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

func (w *capturingResponseWriter) WriteString(s string) (int, error) {
	w.body = append(w.body, []byte(s)...)
	return w.ResponseWriter.WriteString(s)
}

func (w *capturingResponseWriter) Body() []byte {
	return w.body
}
