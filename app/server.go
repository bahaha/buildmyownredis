package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os/signal"
	"strconv"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

var (
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

func NewRedisServer(ctx context.Context, port int) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	logger.Info("Starting TCP server", "port", port)

	l, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(port)))
	if err != nil {
		logger.Error("Failed to bind to port", "error", err)
		os.Exit(1)
	}

	logger.Info("Listening for connections on port", "port", port)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					// Server is shutting down, exit the goroutine
					return
				default:
					logger.Error("Error accepting connection", "error", err)

				}
			}
			go handleConnection(ctx, conn)
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down server")
}

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	logger.Info("Accepted connection", "remote_addr", conn.RemoteAddr())

	select {
	case <-ctx.Done():
		logger.Info("Connection closed due to server shutdown", "remote_addr", conn.RemoteAddr())
		return
	default:
		// Handle the connection
		buf := make([]byte, 1024)

		for {
			command, err := conn.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					logger.Info(
						"Connection closed by remote host",
						"remote_addr",
						conn.RemoteAddr(),
					)
					break
				}
				logger.Error("Error reading from connection", "error", err)
				return
			}

			logger.Info("Received command", "command", command)

			_, err = conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				logger.Error("Error writing to connection", "error", err)
				return
			}
		}
	}
}

func main() {
	ctx := context.Background()
	NewRedisServer(ctx, 6379)
}
