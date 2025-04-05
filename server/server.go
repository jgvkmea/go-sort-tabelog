package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jgvkmea/go-sort-tabelog/interface/controller"
	"github.com/jgvkmea/go-sort-tabelog/interface/controller/middleware/logger"
)

func StartServer(ctx context.Context, addr string, certPath string, keyPath string) {
	log := logger.FromContext(ctx)

	router := newRouter()
	setMiddleware(router)
	setRouting(router)
	s := newServer(router, addr, 30*time.Second)

	idleConnsClosed := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL)
		<-sig

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Infoln("failed to shutdown server: ", err)
		}
		log.Infoln("complete to shutdown")
		close(idleConnsClosed)
	}()

	log.Infoln("start linebot server")
	if err := s.ListenAndServeTLS(certPath, keyPath); err != nil {
		log.Infoln("failed to listen and serve: ", err)
	}
	<-idleConnsClosed
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
