package server

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
	"time"
)

type logEntry struct {
	method     string
	path       string
	remoteAddr string
}

func (l logEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	slog.Info(fmt.Sprintf("%s %s", l.method, l.path), slog.String("remote_addr", l.remoteAddr), slog.Int("status", status), slog.Int("bytes", bytes), slog.Duration("duration", elapsed))
}

func (l logEntry) Panic(v interface{}, stack []byte) {
	slog.Info(fmt.Sprintf("%s %s", l.method, l.path), slog.String("remote_addr", l.remoteAddr), slog.Any("panic", v), slog.String("stack", string(stack)))
}

type logFormatter struct{}

func (l logFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return logEntry{
		method:     r.Method,
		path:       r.URL.Path,
		remoteAddr: r.RemoteAddr,
	}
}
