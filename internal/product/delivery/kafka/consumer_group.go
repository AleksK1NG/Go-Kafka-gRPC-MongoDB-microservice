package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"

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

	writerReadTimeout  = 10 * time.Second
	writerWriteTimeout = 10 * time.Second

	createProductTopic   = "create-product"
	createProductWorkers = 16
	updateProductTopic   = "update-product"
	updateProductWorkers = 16

	deadLetterQueueTopic = "dead-letter-queue"

	productsGroupID = "products_group"
)

// ProductsConsumerGroup struct
type ProductsConsumerGroup struct {
	Brokers    []string
	GroupID    string
	log        logger.Logger
	cfg        *config.Config
	productsUC product.UseCase
}

// NewProductsConsumerGroup constructor
func NewProductsConsumerGroup(
	brokers []string,
	groupID string,
	log logger.Logger,
	cfg *config.Config,
	productsUC product.UseCase,
) *ProductsConsumerGroup {
	return &ProductsConsumerGroup{Brokers: brokers, GroupID: groupID, log: log, cfg: cfg, productsUC: productsUC}
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
			Timeout: dialTimeout,
		},
	})
}

func (pcg *ProductsConsumerGroup) getNewKafkaWriter(topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(pcg.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: -1,
		MaxAttempts:  3,
		Logger:       kafka.LoggerFunc(pcg.log.Infof),
		ErrorLogger:  kafka.LoggerFunc(pcg.log.Errorf),
		Compression:  compress.Snappy,
		ReadTimeout:  writerReadTimeout,
		WriteTimeout: writerWriteTimeout,
	}
	return w
}

func (pcg *ProductsConsumerGroup) consumeCreateProduct(
	ctx context.Context,
	cancel context.CancelFunc,
	groupID string,
	topic string,
	workersNum int,
) {
	r := pcg.getNewKafkaReader(pcg.Brokers, topic, groupID)
	defer cancel()
	defer func() {
		if err := r.Close(); err != nil {
			pcg.log.Errorf("r.Close", err)
			cancel()
		}
	}()

	w := pcg.getNewKafkaWriter(deadLetterQueueTopic)
	defer func() {
		if err := w.Close(); err != nil {
			pcg.log.Errorf("w.Close", err)
			cancel()
		}
	}()

	pcg.log.Infof("Starting consumer group: %v", r.Config().GroupID)

	wg := &sync.WaitGroup{}
	for i := 0; i <= workersNum; i++ {
		wg.Add(1)
		go pcg.createProductWorker(ctx, cancel, r, w, wg, i)
	}
	wg.Wait()
}

func (pcg *ProductsConsumerGroup) consumeUpdateProduct(
	ctx context.Context,
	cancel context.CancelFunc,
	groupID string,
	topic string,
	workersNum int,
) {
	r := pcg.getNewKafkaReader(pcg.Brokers, topic, groupID)
	defer cancel()
	defer func() {
		if err := r.Close(); err != nil {
			pcg.log.Errorf("r.Close", err)
			cancel()
		}
	}()

	w := pcg.getNewKafkaWriter(deadLetterQueueTopic)
	defer func() {
		if err := w.Close(); err != nil {
			pcg.log.Errorf("w.Close", err)
			cancel()
		}
	}()

	pcg.log.Infof("Starting consumer group: %v", r.Config().GroupID)

	wg := &sync.WaitGroup{}
	for i := 0; i <= workersNum; i++ {
		wg.Add(1)
		go pcg.updateProductWorker(ctx, cancel, r, w, wg, i)
	}
	wg.Wait()
}

// RunConsumers run kafka consumers
func (pcg *ProductsConsumerGroup) RunConsumers(ctx context.Context, cancel context.CancelFunc) {
	go pcg.consumeCreateProduct(ctx, cancel, productsGroupID, createProductTopic, createProductWorkers)
	go pcg.consumeUpdateProduct(ctx, cancel, productsGroupID, updateProductTopic, updateProductWorkers)
}
