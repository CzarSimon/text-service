package main

import (
	"net/http"

	"github.com/CzarSimon/text-service/go/pkg/utils/httputil"
	"github.com/CzarSimon/text-service/go/pkg/utils/logger"
	"go.uber.org/zap"
)

var log = logger.GetDefaultLogger("main").Sugar()

func main() {
	cfg := getConfig()
	e := getEnv(cfg)
	defer e.Close()

	server := newServer(e)
	log.Info("Started text-service on port: " + cfg.port)
	err := server.ListenAndServe()
	if err != nil {
		log.Error("Unexpected error stoped server.", zap.Error(err))
	}
}

func newServer(e *env) *http.Server {
	r := httputil.NewRouter(e.checkHealth)

	r.GET("/v1/texts/key/:key", e.getTextByKey)
	r.GET("/v1/texts/group/:groupId", e.getTextGroup)

	return &http.Server{
		Addr:    ":" + e.cfg.port,
		Handler: r,
	}
}
