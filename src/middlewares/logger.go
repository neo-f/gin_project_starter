package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Logger(c *gin.Context) {
	start := time.Now()
	c.Next()
	path := strings.Builder{}
	path.WriteString(c.Request.URL.Path)
	raw := c.Request.URL.RawQuery
	if raw != "" {
		path.WriteString("?")
		path.WriteString(raw)
	}
	var event *zerolog.Event
	code := c.Writer.Status()
	switch {
	case code >= http.StatusOK && code < http.StatusBadRequest:
		event = log.Info()
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		event = log.Warn()
	default:
		event = log.Error()
	}

	size := c.Writer.Size()
	if size != -1 {
		event = event.Str("body_size", humanize.Bytes(uint64(size)))
	}

	event = event.
		Int("code", code).
		Str("client_ip", c.ClientIP()).
		Str("method", c.Request.Method).
		Str("latency", time.Since(start).String()).
		Str("path", path.String())
	if len(c.Errors) != 0 {
		event = event.Interface("err_message", c.Errors.String())
	}
	event.Msg("REQ")
}
