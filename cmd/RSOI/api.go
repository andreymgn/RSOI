package main

import (
	"log"

	"github.com/andreymgn/RSOI/pkg/tracer"
	"github.com/andreymgn/RSOI/services/api"
	comment "github.com/andreymgn/RSOI/services/comment/proto"
	post "github.com/andreymgn/RSOI/services/post/proto"
	poststats "github.com/andreymgn/RSOI/services/poststats/proto"
	user "github.com/andreymgn/RSOI/services/user/proto"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"google.golang.org/grpc"
)

func runAPI(port int, postAddr, commentAddr, postStatsAddr, userAddr, jaegerAddr string) error {
	tracer, err := tracer.NewTracer("api", jaegerAddr)
	if err != nil {
		log.Fatal(err)
	}

	postConn, err := grpc.Dial(postAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
	if err != nil {
		log.Fatal(err)
	}

	defer postConn.Close()
	pc := post.NewPostClient(postConn)

	commentConn, err := grpc.Dial(commentAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
	if err != nil {
		log.Fatal(err)
	}

	defer commentConn.Close()
	cc := comment.NewCommentClient(commentConn)

	postStatsConn, err := grpc.Dial(postStatsAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
	if err != nil {
		log.Fatal(err)
	}

	defer postStatsConn.Close()
	psc := poststats.NewPostStatsClient(postStatsConn)

	userConn, err := grpc.Dial(userAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
	if err != nil {
		log.Fatal(err)
	}

	defer userConn.Close()
	uc := user.NewUserClient(userConn)

	server := api.NewServer(pc, cc, psc, uc, tracer)
	server.Start(port)

	return nil
}
