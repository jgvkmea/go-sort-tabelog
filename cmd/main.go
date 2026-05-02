package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/jgvkmea/go-sort-tabelog/interface/controller/middleware/logger"
	"github.com/jgvkmea/go-sort-tabelog/server"
)

var (
	address = flag.String("addr", "127.0.0.1", "serve ip address")
	port    = flag.String("port", "8080", "serve port")
	tlsCert = flag.String("tls-cert", os.Getenv("TLS_CERT_PATH"), "path to TLS certificate file; HTTPS is enabled when set together with --tls-key")
	tlsKey  = flag.String("tls-key", os.Getenv("TLS_KEY_PATH"), "path to TLS key file; HTTPS is enabled when set together with --tls-cert")
)

func main() {
	flag.Parse()

	ctx := context.Background()
	ctx = logger.WithLogger(ctx)
	log := logger.FromContext(ctx)

	if (*tlsCert == "") != (*tlsKey == "") {
		log.Errorln("--tls-cert and --tls-key must be set together")
		return
	}

	server.StartServer(ctx, fmt.Sprintf("%s:%s", *address, *port), *tlsCert, *tlsKey)
}
