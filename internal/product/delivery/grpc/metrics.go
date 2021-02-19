package grpc

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	incomingMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_incoming_grpc_messages_total",
		Help: "The total number of incoming gRPC messages",
	})
	successMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_success_incoming_grpc_messages_total",
		Help: "The total number of success incoming success gRPC messages",
	})
	errorMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_error_incoming_grpc_message_total",
		Help: "The total number of error incoming success gRPC messages",
	})
)
