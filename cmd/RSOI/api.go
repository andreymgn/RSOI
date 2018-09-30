package main

import (
	"log"

	"github.com/andreymgn/RSOI/services/api"
	comment "github.com/andreymgn/RSOI/services/comment/proto"
	post "github.com/andreymgn/RSOI/services/post/proto"
	poststats "github.com/andreymgn/RSOI/services/poststats/proto"
	"google.golang.org/grpc"
)

func runAPI(port int, postAddress, commentAddress, postStatsAddress string) error {
	postConn, err := grpc.Dial(postAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer postConn.Close()
	pc := post.NewPostClient(postConn)

	commentConn, err := grpc.Dial(commentAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer commentConn.Close()
	cc := comment.NewCommentClient(commentConn)

	postStatsConn, err := grpc.Dial(postStatsAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer postStatsConn.Close()
	psc := poststats.NewPostStatsClient(postStatsConn)

	server := api.NewServer(pc, cc, psc)
	server.Start(port)

	return nil
}
