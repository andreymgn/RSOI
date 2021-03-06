package main

import (
	"log"

	"github.com/andreymgn/RSOI/pkg/tracer"
	"github.com/andreymgn/RSOI/services/user"
)

const (
	UserAppID     = "UserAPI"
	UserAppSecret = "fzFKf3g6QeIdqbP7"
)

func runUser(port int, connString, jaegerAddr, redisAddr, redisPassword string, redisDB int) error {
	tracer, closer, err := tracer.NewTracer("user", jaegerAddr)
	defer closer.Close()
	if err != nil {
		log.Fatal(err)
	}

	knownKeys := map[string]string{UserAppID: UserAppSecret}

	server, err := user.NewServer(connString, redisAddr, redisPassword, redisDB, knownKeys)
	if err != nil {
		log.Fatal(err)
	}

	return server.Start(port, tracer)
}
