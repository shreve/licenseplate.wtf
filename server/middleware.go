package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// Capture the response code whenever it's written so it can be retrieved
type statusCodeLogger struct {
	w    http.ResponseWriter
	Code int
}

func (s statusCodeLogger) Header() http.Header {
	return s.w.Header()
}

func (s statusCodeLogger) Write(content []byte) (int, error) {
	return s.w.Write(content)
}

func (s statusCodeLogger) WriteHeader(statusCode int) {
	s.Code = statusCode
	s.w.WriteHeader(statusCode)
}

var units = []string{"ns", "µs", "ms", "s"}

func printDuration(dur time.Duration) string {
	d := float64(dur)
	u := 0
	for u < len(units) && d > 1000 {
		u += 1
		d /= 1000.0
	}
	return fmt.Sprintf("%0.2f%s", d, units[u])
}

var defaultLogger = log.Default()

type prefixWriter struct {
	baseWriter io.Writer
}

func (p prefixWriter) Write(bytes []byte) (int, error) {
	p.baseWriter.Write([]byte("| "))
	return p.baseWriter.Write(bytes)
}

var requestGroupLogger = prefixWriter{defaultLogger.Writer()}

func init() {
	log.SetFlags(log.Flags() ^ log.Lmsgprefix)
}

var logLock sync.Mutex

// Log details about each request
func (s *server) logging(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := statusCodeLogger{w, 200}
		logLock.Lock()
		defer logLock.Unlock()
		log.Printf("┌-> Started [%s] %s", r.Method, r.URL.Path)
		log.SetPrefix("|   ")
		start := time.Now()
		f.ServeHTTP(rw, r)
		log.SetPrefix("")
		log.Printf("└-> Completed [%s] %d %s (%s)\n\n", r.Method, rw.Code, r.URL.Path, time.Now().Sub(start))
	})
}
