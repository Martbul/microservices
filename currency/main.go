package main

import (
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/martbul/currency/data"
	protos "github.com/martbul/currency/protos/currency"
	"github.com/martbul/currency/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	rates, err := data.NewRates(log)
	if err != nil {
		log.Error("Unable to generate rates", "error", err)
		os.Exit(1)
	}

	gs := grpc.NewServer()
	cs := server.NewCurrency(rates, log)

	protos.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	gs.Serve(l)

}
