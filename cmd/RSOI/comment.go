package main

import (
	"log"

	"github.com/andreymgn/RSOI/pkg/tracer"
	"github.com/andreymgn/RSOI/services/comment"
)

const (
	CommentAppID     = "CommentAPI"
	CommentAppSecret = "PT6RUHLokksaBdIj"
)

func runComment(port int, connString, jaegerAddr, redisAddr, redisPassword string, redisDB int) error {
	tracer, closer, err := tracer.NewTracer("comment", jaegerAddr)
	defer closer.Close()
	if err != nil {
		log.Fatal(err)
	}

	knownKeys := map[string]string{CommentAppID: CommentAppSecret}

	server, err := comment.NewServer(connString, redisAddr, redisPassword, redisDB, knownKeys)
	if err != nil {
		log.Fatal(err)
	}
	return server.Start(port, tracer)
}
