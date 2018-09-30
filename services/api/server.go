package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	comment "github.com/andreymgn/RSOI/services/comment/proto"
	post "github.com/andreymgn/RSOI/services/post/proto"
	poststats "github.com/andreymgn/RSOI/services/poststats/proto"
	"github.com/gorilla/mux"
)

type Server struct {
	router          *mux.Router
	postClient      post.PostClient
	commentClient   comment.CommentClient
	postStatsClient poststats.PostStatsClient
}

// NewServer returns new instance of Server
func NewServer(pc post.PostClient, cc comment.CommentClient, psc poststats.PostStatsClient) *Server {
	return &Server{mux.NewRouter(), pc, cc, psc}
}

// Start starts HTTP server which can shut down gracefully
func (s *Server) Start(port int) {
	s.routes()
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
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
