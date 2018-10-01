package main

import (
	"log"

	"github.com/andreymgn/RSOI/pkg/tracer"
	"github.com/andreymgn/RSOI/services/post"
)

func runPost(port int, connString, jaegerAddr string) error {
	tracer, err := tracer.NewTracer("post", jaegerAddr)
	if err != nil {
		log.Fatal(err)
	}

	server, err := post.NewServer(connString)
	if err != nil {
		log.Fatal(err)
	}
	return server.Start(port, tracer)
}
