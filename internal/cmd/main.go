package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/skamenetskiy/messages/api"
	"github.com/skamenetskiy/messages/internal/database"
	"github.com/skamenetskiy/messages/internal/database/repo"
	"github.com/skamenetskiy/messages/internal/service"
	"google.golang.org/grpc"
)

func main() {
	db, err := database.NewDefault()
	if err != nil {
		log.Fatalln("failed to init database:", err)
	}

	r := repo.New(db)
	s := service.New(r)

	listen(s)
}

func listen(s *service.Service) {
	c := context.Background()
	g := grpc.NewServer()
	m := runtime.NewServeMux()

	api.RegisterMessagesAPIServer(g, s)

	if err := api.RegisterMessagesAPIHandlerServer(c, m, s); err != nil {
		log.Fatalln("failed to register grpc-gateway handler", err)
	}

	h := &http.Server{
		Addr:    getAddr("HTTP", "8080"),
		Handler: m,
	}

	go startGRPC(g)
	go startHTTP(h)

	o := make(chan os.Signal, 1)
	signal.Notify(o, os.Interrupt, os.Kill)
	sig := <-o
	log.Printf("Received signal '%s', shutting down", sig)
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	if err := h.Shutdown(c); err != nil {
		log.Printf("failed to shutdown http server gracefully: %s\n", err)
	}
	g.GracefulStop()
}

func startGRPC(g *grpc.Server) {
	addr := getAddr("GRPC", "50051")
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen grpc on %s: %s\n", addr, err)
	}
	if err = g.Serve(l); err != nil {
		log.Fatalf("failed to serve grpc on %s: %s\n", addr, err)
	}
}

func startHTTP(s *http.Server) {
	if err := s.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("failed to server http:", err)
		}
	}
}

func getAddr(prefix, defaultPort string) (addr string) {
	defer func() { log.Printf("%s listen address is %s\n", prefix, addr) }()
	host := os.Getenv(prefix + "_HOST")
	port := os.Getenv(prefix + "_PORT")
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = defaultPort
	}
	addr = host + ":" + port
	return
}
