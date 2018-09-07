package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type server struct {
	router *http.ServeMux
}

func newServer() *server {
	return &server{http.NewServeMux()}
}

func (s *server) routes() {
	s.router.HandleFunc("/hello/", hello)
	s.router.HandleFunc("/add/", add)
}

func (s *server) run() {
	s.routes()
	port := os.Getenv("PORT")
	var addr string
	if port == "" {
		addr = "localhost:8080"
	} else {
		addr = ":" + port
	}
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Hello, world"))
}

func add(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Lhs int64 `json:"lhs"`
		Rhs int64 `json:"rhs"`
	}

	type response struct {
		Result int64 `json:"result"`
	}

	var req request
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		http.Error(w, "Unprocessable Entity", 422)
	}

	resp := response{req.Lhs + req.Rhs}
	j, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(j)
}
