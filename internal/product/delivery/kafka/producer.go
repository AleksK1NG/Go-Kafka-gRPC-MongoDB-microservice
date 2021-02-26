package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"

	"github.com/AleksK1NG/products-microservice/config"
	"github.com/AleksK1NG/products-microservice/pkg/logger"
)

// ProductsProducer interface
type ProductsProducer interface {
	PublishCreate(ctx context.Context, msgs ...kafka.Message) error
	PublishUpdate(ctx context.Context, msgs ...kafka.Message) error
	Close()
	Run()
	GetNewKafkaWriter(topic string) *kafka.Writer
}

type productsProducer struct {
	log          logger.Logger
	cfg          *config.Config
	createWriter *kafka.Writer
	updateWriter *kafka.Writer
}

// NewProductsProducer constructor
func NewProductsProducer(log logger.Logger, cfg *config.Config) *productsProducer {
	return &productsProducer{log: log, cfg: cfg}
}

// GetNewKafkaWriter Create new kafka writer
func (p *productsProducer) GetNewKafkaWriter(topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(p.cfg.Kafka.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: writerRequiredAcks,
		MaxAttempts:  writerMaxAttempts,
		Logger:       kafka.LoggerFunc(p.log.Debugf),
		ErrorLogger:  kafka.LoggerFunc(p.log.Errorf),
		Compression:  compress.Snappy,
		ReadTimeout:  writerReadTimeout,
		WriteTimeout: writerWriteTimeout,
	}
	return w
}

// Run init producers writers
func (p *productsProducer) Run() {
	p.createWriter = p.GetNewKafkaWriter(createProductTopic)
	p.updateWriter = p.GetNewKafkaWriter(updateProductTopic)
}

// Close close writers
func (p productsProducer) Close() {
	p.createWriter.Close()
	p.updateWriter.Close()
}

// PublishCreate publish messages to create topic
func (p *productsProducer) PublishCreate(ctx context.Context, msgs ...kafka.Message) error {
	return p.createWriter.WriteMessages(ctx, msgs...)
}

// PublishUpdate publish messages to update topic
func (p *productsProducer) PublishUpdate(ctx context.Context, msgs ...kafka.Message) error {
	return p.updateWriter.WriteMessages(ctx, msgs...)
}
