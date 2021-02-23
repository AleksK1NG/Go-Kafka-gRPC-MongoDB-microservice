package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	GRPC_PORT = "GRPC_PORT"
	HTTP_PORT = "HTTP_PORT"
)

// Config of application
type Config struct {
	AppVersion string
	Server     Server
	Logger     Logger
	Jaeger     Jaeger
	Metrics    Metrics
	MongoDB    MongoDB
	Kafka      Kafka
	Http       Http
	Redis      Redis
}

// Server config
type Server struct {
	Port              string
	Development       bool
	Timeout           time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	MaxConnectionIdle time.Duration
	MaxConnectionAge  time.Duration
	Kafka             Kafka
}

type Http struct {
	Port              string
	PprofPort         string
	Timeout           time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	CookieLifeTime    int
	SessionCookieName string
}

// Logger config
type Logger struct {
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// Metrics config
type Metrics struct {
	Port        string
	URL         string
	ServiceName string
}

// Jaeger config
type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

type MongoDB struct {
	URI      string
	User     string
	Password string
	DB       string
}

type Kafka struct {
	Brokers []string
}

type Redis struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultDB string
	MinIdleConn    int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

func exportConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if os.Getenv("MODE") == "DOCKER" {
		viper.SetConfigName("config-docker.yml")
	} else {
		viper.SetConfigName("config.yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// ParseConfig Parse config file
func ParseConfig() (*Config, error) {
	if err := exportConfig(); err != nil {
		return nil, err
	}

	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	gRPCPort := os.Getenv(GRPC_PORT)
	if gRPCPort != "" {
		c.Server.Port = gRPCPort
	}

	httpPort := os.Getenv(HTTP_PORT)
	if httpPort != "" {
		c.Http.Port = httpPort
	}

	return &c, nil
}
