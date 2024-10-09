package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/recorder"

	"github.com/satori/uuid"
	"go.uber.org/zap"
)

type AccessLogStart struct {
	UserAgent     string
	RealIp        string
	ContentLength int64
	URI           string
	Method        string
	StartTime     string
	RequestId     string
}

type AccessLogEnd struct {
	LatencyHuman   string
	LatencyMs      string
	EndTime        string
	RequestId      string
	ResponseStatus int
}

func Access(handler http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.NewV4().String()
		timeNow := time.Now().UTC()
		LogInitRequest(r, logger, timeNow, requestId)

		unauthId, err := functions.GetCookieUnauthIdValue(r)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			return
		}
		ctx := context.WithValue(r.Context(), cnst.UnauthIdCookieName, unauthId)
		ctx = context.WithValue(ctx, cnst.RequestId, requestId)
		r = r.WithContext(ctx)

		rec, ok := w.(*recorder.ResponseWriter)
		if !ok {
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, r)

		LogEndRequest(logger, timeNow, requestId, rec.StatusCode)
	})
}

func LogInitRequest(r *http.Request, logger *zap.Logger, timeNow time.Time, requestId string) {
	msg := fmt.Sprintf("init request %s", requestId)
	startLog := &AccessLogStart{
		UserAgent:     r.UserAgent(),
		RealIp:        r.Header.Get("X-Real-IP"),
		ContentLength: r.ContentLength,
		URI:           r.RequestURI,
		Method:        r.Method,
		StartTime:     timeNow.Format(cnst.Timestamptz),
	}

	logger.Info(msg,
		zap.String("user_agent", startLog.UserAgent),
		zap.String("real_ip", startLog.RealIp),
		zap.Int64("content_length", startLog.ContentLength),
		zap.String("uri", startLog.URI),
		zap.String("method", startLog.Method),
		zap.String("start_time", startLog.StartTime),
		zap.String("request_id", requestId),
	)
}

func LogEndRequest(logger *zap.Logger, timeNow time.Time, requestId string, responseStatus int) {
	msg := fmt.Sprintf("request done %s", requestId)
	endLog := &AccessLogEnd{
		EndTime:        timeNow.Format(cnst.Timestamptz),
		LatencyHuman:   time.Since(timeNow).String(),
		LatencyMs:      time.Since(timeNow).String(),
		ResponseStatus: responseStatus,
	}
	logger.Info(msg,
		zap.String("end_time", endLog.EndTime),
		zap.String("latency_human", endLog.LatencyHuman),
		zap.String("latency_human_ms", endLog.LatencyMs),
		zap.String("request_id", requestId),
		zap.Int("response_status", responseStatus),
	)
}
