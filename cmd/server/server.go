package main

import (
	"flag"
	"log"
	"mafia-grpc/api"
	"mafia-grpc/internal/server"
	"mafia-grpc/internal/session"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func serverOptions(args Args) (opt []grpc.ServerOption) {
	opt = append(opt, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: args.KeepaliveTimeout,
		Time:              args.KeepaliveTimeout,
		Timeout:           args.KeepaliveTimeout,
	}))
	return
}

func initListenerAndServer(args Args) (l net.Listener, s *grpc.Server) {
	l, err := net.Listen("tcp", args.Addr)
	if err != nil {
		log.Fatalf("Failed init listener: %v", err)
	}

	s = grpc.NewServer(serverOptions(args)...)

	api.RegisterMafiaServiceServer(s, server.NewServer(server.ServerConfig{
		GameSessionConfig: session.GameSessionConfig{
			PlayersInSession: args.PlayersInSession,
			OpenVoting:       args.OpenVoting,
		},
	}))
	return
}

type Args struct {
	Addr             string
	KeepaliveTimeout time.Duration
	OpenVoting       bool
	PlayersInSession uint
}

func parseArgs() (args Args) {
	flag.StringVar(&args.Addr, "addr", "0.0.0.0:8080", "address")
	flag.DurationVar(&args.KeepaliveTimeout, "keepalive", time.Second, "keap-alive timeout")
	flag.BoolVar(&args.OpenVoting, "openvoting", false, "Will the vote be open")
	flag.UintVar(&args.PlayersInSession, "ssize", 4, "Number of players in one session (must be at least 4)")
	flag.Parse()
	return
}

func main() {
	args := parseArgs()
	l, s := initListenerAndServer(args)
	defer l.Close()
	s.Serve(l)
}
