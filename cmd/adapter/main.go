// Copyright 2025 SGNL.ai, Inc.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/server"
	"github.com/sgnl-ai/sample-adapter/pkg/scim"

	"google.golang.org/grpc"
)

var (
	// Port is the port at which the gRPC server will listen.
	Port = flag.Int("port", 8080, "The server port")

	// Timeout is the timeout for the HTTP client used to make requests to the datasource (seconds).
	Timeout = flag.Int("timeout", 30, "The timeout for the HTTP client used to make requests to the datasource (seconds)")
)

func main() {
	logger := log.New(os.Stdout, "adapter", log.Lmicroseconds|log.LUTC|log.Lshortfile)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *Port))
	if err != nil {
		logger.Fatalf("Failed to open server port: %v", err)
	}

	timeout := time.Duration(*Timeout) * time.Second

	s := grpc.NewServer()
	stop := make(chan struct{})
	adapterServer := server.New(stop)

	// Register the adapters
	server.RegisterAdapter(
		adapterServer,
		"SCIM2.0-1.0.0",
		scim.NewAdapter(scim.NewClient(&http.Client{
			Timeout: timeout,
		})),
	)

	api_adapter_v1.RegisterAdapterServer(s, adapterServer)

	logger.Printf("Started adapter gRPC server on port %d", *Port)

	if err := s.Serve(listener); err != nil {
		close(stop)

		logger.Fatalf("Failed to listen on server port: %v", err)
	}
}
