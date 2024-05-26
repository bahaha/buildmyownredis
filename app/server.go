package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
)

var (
	logger            = slog.New(slog.NewTextHandler(os.Stdout, nil))
	defaultListenPort = 6379
)

type Config struct {
	Port int
}

type Server struct {
	Config
	listener net.Listener
	peers    map[*Peer]bool

	cmdParser Parser
	storage   Storage
}

func NewRedis(ctx context.Context, cfg Config) *Server {
	if cfg.Port == 0 {
		cfg.Port = defaultListenPort
	}

	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		cmdParser: NewResp(),
		storage:   NewMemoryStorage(),
	}
}

func (s *Server) Start(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	l, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(s.Port)))

	if err != nil {
		slog.Error("Failed to bind to port", "error", err)
		os.Exit(1)
	}

	slog.Info("Listening on port", "addr", l.Addr())

	s.listener = l

	go func() {
		<-ctx.Done()
		slog.Info("Shutting down server")
		s.listener.Close()
		for peer := range s.peers {
			peer.Close()
		}
	}()
	s.Listen(ctx)

	return nil
}

func (s *Server) Listen(ctx context.Context) error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				slog.Error("Error accepting connection", "error", err)
			}
			continue
		}

		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	slog.Info("Accepted connection", "remote_addr", conn.RemoteAddr())

	peer := NewPeer(conn, s.cmdParser, s.storage)
	s.peers[peer] = true

	defer func() {
		delete(s.peers, peer)
	}()

	select {
	case <-ctx.Done():
		slog.Info("Connection closed due to server shutdown", "remote_addr", conn.RemoteAddr())
		return
	default:
		peer.WaitForCommand()
	}
}

var (
	port int
)

func init() {
	flag.IntVar(&port, "port", defaultListenPort, "port to listen on")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	server := NewRedis(ctx, Config{Port: port})
	server.Start(ctx)
}
