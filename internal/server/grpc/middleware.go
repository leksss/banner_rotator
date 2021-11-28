package internalgrpc

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/leksss/banner_rotator/internal/domain/interfaces"
)

func loggingMiddleware(next http.Handler, log interfaces.Log) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := &responseObserver{ResponseWriter: w}

		next.ServeHTTP(o, r)

		addr := r.RemoteAddr
		if i := strings.LastIndex(addr, ":"); i != -1 {
			addr = addr[:i]
		}

		log.Info(
			fmt.Sprintf("%s [%s] %s %s %s %d %d %q %q",
				addr,
				time.Now().Format("02/Jan/2006:15:04:05 -0700"),
				r.Method,
				r.URL,
				r.Proto,
				o.status,
				o.written,
				r.Referer(),
				r.UserAgent(),
			),
		)
	})
}

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}
