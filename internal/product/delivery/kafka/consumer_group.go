package kafka

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"

	"github.com/AleksK1NG/products-microservice/config"
	"github.com/AleksK1NG/products-microservice/internal/models"
	"github.com/AleksK1NG/products-microservice/internal/product"
	"github.com/AleksK1NG/products-microservice/pkg/logger"
)

// ProductsConsumerGroup struct
type ProductsConsumerGroup struct {
	Brokers    []string
	GroupID    string
	log        logger.Logger
	cfg        *config.Config
	productsUC product.UseCase
	validate   *validator.Validate
}

// NewProductsConsumerGroup constructor
func NewProductsConsumerGroup(
	brokers []string,
	groupID string,
	log logger.Logger,
	cfg *config.Config,
	productsUC product.UseCase,
	validate *validator.Validate,
) *ProductsConsumerGroup {
	return &ProductsConsumerGroup{
		Brokers:    brokers,
		GroupID:    groupID,
		log:        log,
		cfg:        cfg,
		productsUC: productsUC,
		validate:   validate,
	}
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
		Logger:                 kafka.LoggerFunc(pcg.log.Debugf),
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
		RequiredAcks: writerRequiredAcks,
		MaxAttempts:  writerMaxAttempts,
		Logger:       kafka.LoggerFunc(pcg.log.Debugf),
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

func (pcg *ProductsConsumerGroup) publishErrorMessage(ctx context.Context, w *kafka.Writer, m kafka.Message, err error) error {
	errMsg := &models.ErrorMessage{
		Offset:    m.Offset,
		Error:     err.Error(),
		Time:      m.Time.UTC(),
		Partition: m.Partition,
		Topic:     m.Topic,
	}

	errMsgBytes, err := json.Marshal(errMsg)
	if err != nil {
		return err
	}

	return w.WriteMessages(ctx, kafka.Message{
		Value: errMsgBytes,
	})
}

// RunConsumers run kafka consumers
func (pcg *ProductsConsumerGroup) RunConsumers(ctx context.Context, cancel context.CancelFunc) {
	go pcg.consumeCreateProduct(ctx, cancel, productsGroupID, createProductTopic, createProductWorkers)
	go pcg.consumeUpdateProduct(ctx, cancel, productsGroupID, updateProductTopic, updateProductWorkers)
}
