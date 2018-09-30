package main

import (
	"log"

	"github.com/andreymgn/RSOI/services/post"
)

func runPost(port int, connString string) error {
	server, err := post.NewServer(connString)
	if err != nil {
		log.Fatal(err)
	}
	return server.Start(port)
}
