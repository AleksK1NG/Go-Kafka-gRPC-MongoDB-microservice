package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/AleksK1NG/products-microservice/config"
	"github.com/AleksK1NG/products-microservice/internal/product"
	"github.com/AleksK1NG/products-microservice/pkg/logger"
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
)

// ProductsConsumerGroup
type ProductsConsumerGroup struct {
	Brokers    []string
	GroupID    string
	Topic      string
	log        logger.Logger
	cfg        config.Config
	productsUC product.UseCase
}

// NewProductsConsumerGroup constructor
func NewProductsConsumerGroup(
	brokers []string,
	groupID string,
	topic string,
	log logger.Logger,
	cfg config.Config,
	productsUC product.UseCase,
) *ProductsConsumerGroup {
	return &ProductsConsumerGroup{Brokers: brokers, GroupID: groupID, Topic: topic, log: log, cfg: cfg, productsUC: productsUC}
}

func (pcg *ProductsConsumerGroup) getNewKafkaReader(kafkaURL []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                kafkaURL,
		GroupID:                groupID,
		Topic:                  topic,
		MinBytes:               minBytes,
		MaxBytes:               maxBytes,
		QueueCapacity:          queueCapacity,
		HeartbeatInterval:      heartbeatInterval,
		CommitInterval:         commitInterval,
		PartitionWatchInterval: partitionWatchInterval,
		Logger:                 kafka.LoggerFunc(pcg.log.Infof),
		ErrorLogger:            kafka.LoggerFunc(pcg.log.Errorf),
		MaxAttempts:            maxAttempts,
		Dialer: &kafka.Dialer{
			Timeout:         dialTimeout,
			LocalAddr:       nil,
			DualStack:       false,
			FallbackDelay:   0,
			KeepAlive:       0,
			Resolver:        nil,
			TLS:             nil,
			SASLMechanism:   nil,
			TransactionalID: "",
		},
	})
}

func (pcg *ProductsConsumerGroup) consumeCreateProduct(ctx context.Context, cancel context.CancelFunc, topic string, workersNum int) {
	r := pcg.getNewKafkaReader(pcg.Brokers, topic, pcg.GroupID)
	defer cancel()
	defer func() {
		if err := r.Close(); err != nil {
			pcg.log.Errorf("r.Close", err)
			cancel()
		}
	}()

	pcg.log.Infof("Starting consumer group: %v", r.Config().GroupID)

	wg := &sync.WaitGroup{}
	for i := 0; i <= workersNum; i++ {
		wg.Add(1)
		go pcg.createProductWorker(ctx, r, wg, i)
	}
	wg.Wait()
}

func (pcg *ProductsConsumerGroup) consumeUpdateProduct(ctx context.Context, cancel context.CancelFunc, topic string, workersNum int) {
	r := pcg.getNewKafkaReader(pcg.Brokers, topic, pcg.GroupID)
	defer cancel()
	defer func() {
		if err := r.Close(); err != nil {
			pcg.log.Errorf("r.Close", err)
			cancel()
		}
	}()

	pcg.log.Infof("Starting consumer group: %v", r.Config().GroupID)

	wg := &sync.WaitGroup{}
	for i := 0; i <= workersNum; i++ {
		wg.Add(1)
		go pcg.updateProductWorker(ctx, r, wg, i)
	}
	wg.Wait()
}

func (pcg *ProductsConsumerGroup) RunConsumers(ctx context.Context, cancel context.CancelFunc) {
	go pcg.consumeCreateProduct(ctx, cancel, "create_product", 16)
	go pcg.consumeUpdateProduct(ctx, cancel, "update_product", 16)
}
