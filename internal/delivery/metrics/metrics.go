package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	TotalNumberOfRequests     prometheus.Counter
	RequestTime 			  *prometheus.HistogramVec
	NumberOfSpecificRequests  *prometheus.CounterVec
	MicroserviceTimeout		  *prometheus.HistogramVec
	DatabaseDuration		  *prometheus.HistogramVec
}	

const namespace = "gateway"

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		// business metrics 
		// total number of hits
		TotalNumberOfRequests: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name: "number_of_requests",
				Help: "number of requests",
			},
		),
		// number of requests divided by request method and response status
		NumberOfSpecificRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name: "number_of_specific_request",
				Help: "number of request; you can find the number by request method and response status",
			},
			[]string{"method", "path", "status"},
		),
		// request time can be filtered by request method, url path and response status
		RequestTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:    "time_of_request",
				Help:    "HTTP request  duration in milliseconds",
				Buckets: []float64{10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
			},
			[]string{"method", "path", "status"},
		),
		// request time can be filtered by request method, url path and response status
		MicroserviceTimeout: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:    "time_of_request",
				Help:    "gRPC request duration in milliseconds",
				Buckets: []float64{10, 25, 50, 100, 250, 500, 1000},
			},
			[]string{"microservice"},
		),
		DatabaseDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:    "database_duration_ms",
				Help:    "Database request duration in milliseconds",
				Buckets: []float64{10, 25, 50, 100, 250, 500, 1000},
			},
			[]string{"operation"},
		),
	}
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(m.TotalNumberOfRequests, m.NumberOfSpecificRequests, m.RequestTime)
	return m
}

// что надо мониторить
// ! системные метрики
// - ЦПУ сервера (sys-user-idle)
// - память сервера (свободно - занято)
// - место на диске (всего - занято - заводно)
// - трафик (вход - выход)
// ! инфраструктура
// - цпу приложения (sys-user)
// - память приложения (rss, pss)
// - если дятнетесь - сисколлы (read write)
// ! бизнесовые метрики
// - хиты +
// - хиты с разделением по кодам ответов + хендлерам + путям +
// - тайминги ответов по хендлерам +
// - количество горутин + 
// - размер хипа (heap, динамическая память) +  
// - тайминги ответов внешних систем (микросервисы, базы и все такое) 
// - кол-во ошибок внешних систем
// - разлиные бизнесовые операции (почти хиты, н оне по ручкам)
