package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	gologger "github.com/nrf24l01/go-logger"
)

func RequestLogger(l *gologger.Logger) echo.MiddlewareFunc {
	return echoMw.RequestLoggerWithConfig(echoMw.RequestLoggerConfig{
		LogLatency:   true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogURI:       true,
		LogUserAgent: true,
		LogStatus:    true,
		LogError:     true,
		HandleError:  true,
		LogValuesFunc: func(c echo.Context, v echoMw.RequestLoggerValues) error {
			traceId := ""
			if tid := c.Get("traceId"); tid != nil {
				if s, ok := tid.(string); ok {
					traceId = s
				} else {
					traceId = fmt.Sprint(tid)
				}
			}

			// Base colors
			reset := gologger.Reset
			methodColor := gologger.BrightBlue
			pathColor := gologger.Dim
			traceColor := gologger.Cyan

			// Status on code
			statusColor := gologger.BrightGreen
			switch {
			case v.Status >= 500:
				statusColor = gologger.BrightRed
			case v.Status >= 400:
				statusColor = gologger.BrightYellow
			case v.Status >= 300:
				statusColor = gologger.BrightBlue
			default:
				statusColor = gologger.BrightGreen
			}

			// Combine colors and fields
			methodField := fmt.Sprintf("%s%s%s", methodColor, v.Method, reset)
			statusField := fmt.Sprintf("%s%d%s", statusColor, v.Status, reset)
			pathField := fmt.Sprintf("%s%s%s", pathColor, v.URI, reset)
			latencyField := fmt.Sprintf("%s", v.Latency)

			tracePart := ""
			if traceId != "" {
				tracePart = fmt.Sprintf("(%s%s%s)", traceColor, traceId, reset)
			}

			// Add req id if exists
			reqIDPart := ""
			if v.RequestID != "" {
				reqIDPart = fmt.Sprintf(" request_id=%s", v.RequestID)
			}

			// Build final message
			msg := fmt.Sprintf("%s %s %s %s %s%s",
				methodField, statusField, pathField, latencyField, tracePart, reqIDPart)

			// Log message
			httpType := gologger.LogType("HTTP")
			if v.Error == nil {
				l.Log(gologger.LevelInfo, httpType, msg, v.RequestID)
			} else {
				full := fmt.Sprintf("%s error=%s", msg, v.Error.Error())
				l.Log(gologger.LevelError, httpType, full, v.RequestID)
			}

			return nil
		},
	})
}
