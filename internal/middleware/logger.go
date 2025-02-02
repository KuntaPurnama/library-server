package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"time"
)

var filePathName = "logs/library-backend.log"

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	rw.body.Write(p)
	return rw.ResponseWriter.Write(p)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		startTime := time.Now()

		ip := c.ClientIP()
		endpoint := c.Request.RequestURI
		method := c.Request.Method
		headers := c.Request.Header

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		logEntry := fmt.Sprintf(
			"[%s] %s %s %s Request ID: %s\nHeaders: %v\nBody: %s\n",
			time.Now().Format(time.RFC3339),
			method,
			ip,
			endpoint,
			requestID,
			headers,
			string(bodyBytes),
		)

		if err := os.MkdirAll("logs", os.ModePerm); err != nil {
			log.Fatal("Error creating logs directory: ", err)
		}

		logFile, err := os.OpenFile(filePathName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Error opening log file: ", err)
		}
		defer logFile.Close()

		if _, err := logFile.WriteString(logEntry); err != nil {
			log.Fatal("Error writing to log file: ", err)
		}

		bodyBuffer := bytes.NewBufferString("")
		rw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bodyBuffer,
		}
		c.Writer = rw

		c.Next()

		duration := time.Since(startTime)
		logResponse := fmt.Sprintf(
			"[%s] %s %s %s Request ID: %s\nDuration: %v\nResponse Body: %s\n",
			time.Now().Format(time.RFC3339),
			method,
			ip,
			endpoint,
			requestID,
			duration,
			rw.body.String(),
		)

		logFile, err = os.OpenFile(filePathName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Error opening log file: ", err)
		}
		defer logFile.Close()

		if _, err := logFile.WriteString(logResponse); err != nil {
			log.Fatal("Error writing to log file: ", err)
		}
	}
}
