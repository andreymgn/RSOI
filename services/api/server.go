package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/andreymgn/RSOI/pkg/tracer"
	comment "github.com/andreymgn/RSOI/services/comment/proto"
	post "github.com/andreymgn/RSOI/services/post/proto"
	poststats "github.com/andreymgn/RSOI/services/poststats/proto"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/cors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	router          *tracer.TracedRouter
	postClient      post.PostClient
	commentClient   comment.CommentClient
	postStatsClient poststats.PostStatsClient
}

// NewServer returns new instance of Server
func NewServer(pc post.PostClient, cc comment.CommentClient, psc poststats.PostStatsClient, tr opentracing.Tracer) *Server {
	return &Server{tracer.NewRouter(tr), pc, cc, psc}
}

func handleRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch st.Code() {
	case codes.NotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	case codes.InvalidArgument:
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func setContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// Start starts HTTP server which can shut down gracefully
func (s *Server) Start(port int) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                                         // All origins
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}, // Allowing only get, just an example
		AllowedHeaders:   []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	})
	s.router.Mux.Use(setContentType)
	s.routes()
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      c.Handler(s.router),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt)
	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
