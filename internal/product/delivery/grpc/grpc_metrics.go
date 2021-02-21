package grpc

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	successMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_success_incoming_grpc_messages_total",
		Help: "The total number of success incoming success gRPC messages",
	})
	errorMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_error_incoming_grpc_message_total",
		Help: "The total number of error incoming success gRPC messages",
	})
	createMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_create_incoming_grpc_requests_total",
		Help: "The total number of incoming create product gRPC messages",
	})
	updateMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_update_incoming_grpc_requests_total",
		Help: "The total number of incoming update product gRPC messages",
	})
	getByIdMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_get_by_id_incoming_grpc_requests_total",
		Help: "The total number of incoming get by id product gRPC messages",
	})
	searchMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_search_incoming_grpc_requests_total",
		Help: "The total number of incoming search products gRPC messages",
	})
)
