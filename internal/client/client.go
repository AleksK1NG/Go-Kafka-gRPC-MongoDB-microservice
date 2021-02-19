package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/AleksK1NG/products-microservice/internal/models"
)

func main() {
	conn, err := kafka.DialContext(context.Background(), "tcp", "localhost:9092")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	w := &kafka.Writer{
		Addr:         kafka.TCP("localhost:9091", "localhost:9092", "localhost:9093"),
		Topic:        "create-product",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: -1,
		MaxAttempts:  3,
	}
	defer w.Close()

	prod := &models.Product{
		CategoryID:  primitive.NewObjectID(),
		Name:        "Book Microservices",
		Description: "Book Microservices Architecture",
		Price:       10.50,
		ImageURL:    nil,
		Photos:      nil,
		Quantity:    1000,
		Rating:      10,
	}

	prodBytes, err := json.Marshal(&prod)
	if err != nil {
		log.Fatal(err)
	}

	msg := kafka.Message{
		Value: prodBytes,
	}
	err = w.WriteMessages(context.Background(), msg)
	if err != nil {
		fmt.Println(err)
	}
}
