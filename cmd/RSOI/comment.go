package main

import (
	"log"

	"github.com/andreymgn/RSOI/pkg/tracer"
	"github.com/andreymgn/RSOI/services/comment"
)

func runComment(port int, connString, jaegerAddr string) error {
	tracer, err := tracer.NewTracer("comment", jaegerAddr)
	if err != nil {
		log.Fatal(err)
	}

	server, err := comment.NewServer(connString)
	if err != nil {
		log.Fatal(err)
	}
	return server.Start(port, tracer)
}
