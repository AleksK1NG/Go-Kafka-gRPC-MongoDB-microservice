package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AleksK1NG/products-microservice/config"
	"github.com/AleksK1NG/products-microservice/pkg/logger"
)

func main() {
	log.Println("Starting products service")
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Info("Starting user server")
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, DevelopmentMode: %s",
		cfg.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Development,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.AppVersion)

	http.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST: %v", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "ok"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Fatal(http.ListenAndServe(cfg.Server.Port, nil))
}
