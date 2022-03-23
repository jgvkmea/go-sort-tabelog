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
)

func main() {
	ctx := context.Background()
	ctx = logger.WithLogger(ctx)
	log := logger.FromContext(ctx)

	certPath, keyPath := getSSHFilePath()
	if certPath == "" || keyPath == "" {
		log.Errorln("require parameter certPath and keyPath")
		return
	}

	flag.Parse()

	if err := server.StartServer(ctx, fmt.Sprintf("%s:%s", *address, *port), certPath, keyPath); err != nil {
		log.Errorln("failed to start server: ", err)
	}
}

func getSSHFilePath() (certPath string, keyPath string) {
	certPath = os.Getenv("SSH_CERT_PATH")
	keyPath = os.Getenv("SSH_KEY_PATH")
	return
}
