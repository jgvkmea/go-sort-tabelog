package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jgvkmea/go-sort-tabelog/middleware/logger"
	"github.com/jgvkmea/go-sort-tabelog/service"
)

func StartServer(ctx context.Context, address string, port string, certPath string, keyPath string) error {
	log := logger.FromContext(ctx)

	router := mux.NewRouter()
	routing(router)

	s := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%s", address, port),
		WriteTimeout: 30 * time.Second,
	}

	log.Infoln("start linebot server")
	if err := s.ListenAndServeTLS(certPath, keyPath); err != nil {
		log.Errorln("failed to start server")
		return err
	}
	return nil
}

func routing(r *mux.Router) {
	r.Use(logger.Middleware)
	r.HandleFunc("/linebot/tabelog", service.TabelogSearchHandler).Methods("POST")
}
