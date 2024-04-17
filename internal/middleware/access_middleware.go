package middleware

import (
	"fmt"
	"net/http"
	"time"

	cnst "2024_1_kayros/internal/utils/constants"
	"github.com/satori/uuid"
	"go.uber.org/zap"
)

type AccessLogStart struct {
	UserAgent     string
	Host          string
	RealIp        string
	ContentLength int64
	URI           string
	Method        string
	StartTime     string
	RequestId     string
}

type AccessLogEnd struct {
	UserAgent    string
	Host         string
	RealIp       string
	ResponseSize int64
	URI          string
	Method       string
	LatencyHuman string
	LatencyMs    string
	EndTime      string
	RequestId    string
}

func AccessMiddleware(handler http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.NewV4().String()
		timeNow := time.Now().UTC()

		handler.ServeHTTP(w, r)
	})
}

func LogInitRequest(r *http.Request, logger *zap.Logger, timeNow time.Time, requestId string) {
	msg := fmt.Sprintf("init request %s", requestId)
	startLog := &AccessLogStart{
		UserAgent:     r.UserAgent(),
		Host:          r.Host,
		RealIp:        r.Header.Get("X-Real-IP"),
		ContentLength: r.ContentLength,
		URI:           r.RequestURI,
		Method:        r.Method,
		StartTime:     timeNow.Format(cnst.Timestamptz),
	}

	logger.Info(msg,
		zap.String("user_agent", startLog.UserAgent),
		zap.String("host", startLog.Host),
		zap.String("real_ip", startLog.RealIp),
		zap.Int64("content_length", startLog.ContentLength),
		zap.String("uri", startLog.URI),
		zap.String("method", startLog.Method),
		zap.String("start_time", startLog.StartTime),
		zap.String("request_id", requestId),
	)
}

func LogEndRequest(r *http.Request, logger *zap.Logger, timeNow time.Time, requestId string) {
	msg := fmt.Sprintf("request done %s", requestId)
	endLog := &AccessLogEnd{
		UserAgent:    r.UserAgent(),
		Host:         r.Host,
		RealIp:       r.Header.Get("X-Real-IP"),
		URI:          r.RequestURI,
		Method:       r.Method,
		EndTime:      timeNow.Format(cnst.Timestamptz),
		LatencyHuman: time.Since(timeNow).Seconds(),
		LatencyMs, string
		EndTime, string,
	}
	logger.Info(msg,
		zap.String("user_agent", startLog.UserAgent),
		zap.String("host", startLog.Host),
		zap.String("real_ip", startLog.RealIp),
		zap.Int64("content_length", startLog.ContentLength),
		zap.String("uri", startLog.URI),
		zap.String("method", startLog.Method),
		zap.String("remote_ip", startLog.RemoteIp),
		zap.String("start_time", startLog.StartTime),
		zap.String("response_size", star)
	zap.String("request_id", requestId),
)
}
