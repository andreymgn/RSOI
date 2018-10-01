package main

import (
	"log"

	"github.com/andreymgn/RSOI/pkg/tracer"
	"github.com/andreymgn/RSOI/services/poststats"
)

func runPostStats(port int, connString, jaegerAddr string) error {
	tracer, err := tracer.NewTracer("poststats", jaegerAddr)
	if err != nil {
		log.Fatal(err)
	}

	server, err := poststats.NewServer(connString)
	if err != nil {
		log.Fatal(err)
	}
	return server.Start(port, tracer)
}
