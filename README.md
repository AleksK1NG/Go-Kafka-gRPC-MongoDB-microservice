### Golang Kafka gRPC MongoDB microservice example 👋

#### 👨‍💻 Full list what has been used:
* [Kafka](https://github.com/segmentio/kafka-go) - Kafka library in Go
* [gRPC](https://grpc.io/) - gRPC
* [echo](https://github.com/labstack/echo) - Web framework
* [viper](https://github.com/spf13/viper) - Go configuration with fangs
* [go-redis](https://github.com/go-redis/redis) - Type-safe Redis client for Golang
* [zap](https://github.com/uber-go/zap) - Logger
* [validator](https://github.com/go-playground/validator) - Go Struct and Field validation
* [swag](https://github.com/swaggo/swag) - Swagger
* [CompileDaemon](https://github.com/githubnemo/CompileDaemon) - Compile daemon for Go
* [Docker](https://www.docker.com/) - Docker
* [Prometheus](https://prometheus.io/) - Prometheus
* [Grafana](https://grafana.com/) - Grafana
* [Jaeger](https://www.jaegertracing.io/) - Jaeger tracing
* [MongoDB](https://github.com/mongodb/mongo-go-driver) - The Go driver for MongoDB
* [retry-go](https://github.com/avast/retry-go) - Simple golang library for retry mechanism
* [kafdrop](https://github.com/obsidiandynamics/kafdrop) - Kafka Web UI


### Jaeger UI:

http://localhost:16686

### Prometheus UI:

http://localhost:9090

### Grafana UI:

http://localhost:3000

### Kafka UI:

http://localhost:9000/

### Swagger UI:

https://localhost:5007/swagger/index.html

For local development:
```
make local // runs docker-compose.local.yml
make crate_topics // create kafka topics
make mongo // load js init script to mongo docker container
make cert // generate local SLL certificates
make swagger // generate swagger documentation
```