package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	TotalNumberOfRequests    prometheus.Counter
	RequestTime              *prometheus.HistogramVec
	NumberOfSpecificRequests *prometheus.CounterVec
	MicroserviceTimeout      *prometheus.HistogramVec
	MicroserviceErrors       *prometheus.CounterVec
	DatabaseDuration         *prometheus.HistogramVec
	PopularRestaurant        *prometheus.CounterVec
	PopularFood              *prometheus.CounterVec
	PopularCategory          *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer, namespace string) *Metrics {
	m := &Metrics{
		// business metrics
		// total number of hits
		TotalNumberOfRequests: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "number_of_requests",
				Help:      "number of requests",
			},
		),
		// number of requests divided by request method and response status
		NumberOfSpecificRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "number_of_specific_request",
				Help:      "number of request; you can find the number by request method and response status",
			},
			[]string{"method", "path", "status"},
		),
		// request time can be filtered by request method, url path and response status
		RequestTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "time_of_request",
				Help:      "HTTP request  duration in milliseconds",
				Buckets:   []float64{10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
			},
			[]string{"method", "path", "status"},
		),
		// request time can be filtered by request method, url path and response status
		MicroserviceErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "count_error",
				Help:      "count of errors received from microservice",
			},
			[]string{"microservice", "status_code"},
		),
		MicroserviceTimeout: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "time_of_request_microservice",
				Help:      "gRPC request duration in milliseconds",
				Buckets:   []float64{10, 25, 50, 100, 250, 500, 1000, 2500},
			},
			[]string{"microservice"},
		),
		DatabaseDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "database_monolith_duration_ms",
				Help:      "Database request duration in milliseconds",
				Buckets:   []float64{10, 25, 50, 100, 250, 500, 1000},
			},
			[]string{"operation"},
		),
		PopularRestaurant: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "business",
				Name:      "rest_visit",
				Help:      "number of times users visited the restaurant page",
			},
			[]string{"restaurant_id"},
		),
		PopularFood: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "business",
				Name:      "popular_food",
				Help:      "number of times users bought the food",
			},
			[]string{"food_id"},
		),
		PopularCategory: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "business",
				Name:      "popular_category",
				Help:      "number of times users opened the current category",
			},
			[]string{"category_id"},
		),
	}
	reg.MustRegister(getCollectors(m)...)
	return m
}

func getCollectors(m *Metrics) []prometheus.Collector {
	cltrs := make([]prometheus.Collector, 0, 10)
	cltrs = append(cltrs, collectors.NewGoCollector())
	cltrs = append(cltrs, m.TotalNumberOfRequests)
	cltrs = append(cltrs, m.NumberOfSpecificRequests)
	cltrs = append(cltrs, m.RequestTime)
	cltrs = append(cltrs, m.MicroserviceErrors)
	cltrs = append(cltrs, m.MicroserviceTimeout)
	cltrs = append(cltrs, m.DatabaseDuration)
	cltrs = append(cltrs, m.PopularRestaurant)
	cltrs = append(cltrs, m.PopularFood)
	cltrs = append(cltrs, m.PopularCategory)
	return cltrs
}

// что надо мониторить
// ! системные метрики
// - ЦПУ сервера (sys-user-idle)
// - память сервера (свободно - занято)
// - место на диске (всего - занято - заводно)
// - трафик (вход - выход)
// ! инфраструктура
// - цпу приложения (sys-user)
// - память приложения (rss, pss) ?
// - если дятнетесь - сисколлы (read write) +
// ! бизнесовые метрики
// - хиты +
// - хиты с разделением по кодам ответов + хендлерам + путям +
// - тайминги ответов по хендлерам +
// - количество горутин +
// - размер хипа (heap, динамическая память) +
// - тайминги ответов внешних систем (микросервисы, базы и все такое) +
// - кол-во ошибок внешних систем +
// - разлиные бизнесовые операции (почти хиты, н оне по ручкам) +
