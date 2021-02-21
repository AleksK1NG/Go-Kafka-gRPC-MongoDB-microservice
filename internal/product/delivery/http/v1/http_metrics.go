package v1

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	successRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_products_success_incoming_messages_total",
		Help: "The total number of success incoming success HTTP requests",
	})
	errorRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_products_error_incoming_message_total",
		Help: "The total number of error incoming success HTTP requests",
	})
	createRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_products_create_incoming_requests_total",
		Help: "The total number of incoming create product HTTP requests",
	})
	updateRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_products_update_incoming_requests_total",
		Help: "The total number of incoming update product HTTP requests",
	})
	getByIdRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_products_get_by_id_incoming_requests_total",
		Help: "The total number of incoming get by id product HTTP requests",
	})
	searchRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_products_search_incoming_requests_total",
		Help: "The total number of incoming search products HTTP requests",
	})
)
