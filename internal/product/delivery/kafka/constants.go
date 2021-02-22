package kafka

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	incomingMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_incoming_kafka_messages_total",
		Help: "The total number of incoming Kafka messages",
	})
	successMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_success_incoming_kafka_messages_total",
		Help: "The total number of success incoming success Kafka messages",
	})
	errorMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_error_incoming_kafka_message_total",
		Help: "The total number of error incoming success Kafka messages",
	})
)

const (
	minBytes               = 10e3 // 10KB
	maxBytes               = 10e6 // 10MB
	queueCapacity          = 100
	heartbeatInterval      = 3 * time.Second
	commitInterval         = 0
	partitionWatchInterval = 5 * time.Second
	maxAttempts            = 3
	dialTimeout            = 3 * time.Minute

	writerReadTimeout  = 10 * time.Second
	writerWriteTimeout = 10 * time.Second

	createProductTopic   = "create-product"
	createProductWorkers = 3
	updateProductTopic   = "update-product"
	updateProductWorkers = 3

	deadLetterQueueTopic = "dead-letter-queue"

	productsGroupID = "products_group"
)
