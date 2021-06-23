package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"test.com/cache-server/handler"
)

// TODO, config those settings
const (
	writeTimeout = 15 * time.Second
	readTimeout  = 15 * time.Second
	idleTimeout  = 60 * time.Second

	shutdownTimeout = 15 * time.Second
)

// Server implements http server
type Server struct {
	*http.Server

	// TODO include configs
}

// NewServer creates instance of Server
func NewServer(addr string) *Server {
	h, err := handler.NewHandler()
	if err != nil {
		panic("error creating hanlder: " + err.Error())
	}

	srv := &http.Server{
		Handler:      h,
		Addr:         addr,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
	}

	return &Server{srv}
}

// Start starts the Server
func (srv *Server) Start() {
	ch := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		log.Infof("stopping server on %s", srv.Addr)

		// We received a stop signal, shut down.
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown: %v", err)
		}
		close(ch)
	}()

	log.Infof("starting server on %s", srv.Addr)

	// TODO, impelments TLS certs config for https
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("http server error: %v", err)
	}

	<-ch
}

func main() {
	addr := flag.String("addr", ":8080", "cache server addr")
	flag.Parse()

	server := NewServer(*addr)
	server.Start()
}
