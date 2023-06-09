package main

import (
	"flag"
	"log"
	"mafia-grpc/internal/bot"
	"mafia-grpc/internal/cli"
	"mafia-grpc/internal/client"
	"mafia-grpc/internal/event"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func dialOptions(args Args) (opt []grpc.DialOption) {
	//opt = append(opt, grpc.WithInsecure())
	opt = append(opt, grpc.WithBlock())
	opt = append(opt, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return
}

func clientOptions(args Args) (opt client.ClientOptions) {
	return client.ClientOptions{
		Addr:        args.Addr,
		GrpcOptions: dialOptions(args),
		EventHandlers: []event.EventHandler{
			cli.NewEventPrinter(),
		},
	}
}

type Args struct {
	Addr   string
	Auto   bool
	BotLag time.Duration
}

func parseArgs() (args Args) {
	flag.StringVar(&args.Addr, "addr", "0.0.0.0:8080", "address")
	flag.BoolVar(&args.Auto, "auto", false, "automatic selection of actions on each turn")
	flag.DurationVar(&args.BotLag, "lag", time.Millisecond*100, "delay between actions in automatic mode")
	flag.Parse()
	return
}

func main() {
	args := parseArgs()
	if args.Auto {
		bot := bot.NewBot(args.BotLag)
		err := bot.StartClient(clientOptions(args))
		if err != nil {
			log.Fatal(err)
		}
		bot.Serve()
	} else {
		cli := cli.NewCLI()
		err := cli.StartClient(clientOptions(args))
		if err != nil {
			log.Fatal(err)
		}
		cli.Serve()
	}
}
