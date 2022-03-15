package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jgvkmea/go-sort-tabelog/interface/controller"
	"github.com/jgvkmea/go-sort-tabelog/interface/controller/middleware/logger"
)

func StartServer(ctx context.Context, addr string, certPath string, keyPath string) error {
	log := logger.FromContext(ctx)

	router := newRouter()
	setMiddleware(router)
	setRouting(router)

	s := newServer(router, addr, 30*time.Second)

	log.Infoln("start linebot server")
	if err := s.ListenAndServeTLS(certPath, keyPath); err != nil {
		log.Errorln("failed to start server")
		return err
	}
	return nil
}

func newRouter() *mux.Router {
	return mux.NewRouter()
}

func setMiddleware(r *mux.Router) {
	r.Use(logger.Middleware)
}

func setRouting(r *mux.Router) {
	r.HandleFunc("/linebot/tabelog", controller.TabelogSearchHandler).Methods("POST")
}

func newServer(r *mux.Router, addr string, timeout time.Duration) *http.Server {
	return &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: timeout,
	}
}
