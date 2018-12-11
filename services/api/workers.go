package api

import (
	"context"
	"log"
	"time"

	comment "github.com/andreymgn/RSOI/services/comment/proto"
	post "github.com/andreymgn/RSOI/services/post/proto"
	poststats "github.com/andreymgn/RSOI/services/poststats/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) deletePostWorker() {
	ctx := context.Background()
	for {
		uid := <-s.deletePostChannel

		log.Printf("DeletePost rabotyaga: got %s", uid)

		go func() {
			_, err := s.postClient.client.DeletePost(ctx,
				&post.DeletePostRequest{Token: s.postClient.token, Uid: uid},
			)
			if err != nil {
				if st, ok := status.FromError(err); ok {
					switch st.Code() {
					case codes.Unauthenticated:
						err := s.updatePostToken()
						if err != nil {
							panic(err)
						}
						_, err = s.postStatsClient.client.DeletePostStats(ctx,
							&poststats.DeletePostStatsRequest{Token: s.postStatsClient.token, PostUid: uid},
						)
						if err != nil {
							panic(err)
						}
					case codes.Unavailable:
						time.Sleep(time.Second * 5)
						s.deletePostChannel <- uid
					}
				}
			}
		}()
	}
}

func (s *Server) deletePostStatsWorker() {
	ctx := context.Background()
	for {
		uid := <-s.deletePostStatsChannel

		log.Printf("DeletePostStats rabotyaga: got %s", uid)

		go func() {
			_, err := s.postStatsClient.client.DeletePostStats(ctx,
				&poststats.DeletePostStatsRequest{Token: s.postStatsClient.token, PostUid: uid},
			)
			if err != nil {
				if st, ok := status.FromError(err); ok {
					switch st.Code() {
					case codes.Unauthenticated:
						err := s.updatePostStatsToken()
						if err != nil {
							panic(err)
						}
						_, err = s.postStatsClient.client.DeletePostStats(ctx,
							&poststats.DeletePostStatsRequest{Token: s.postStatsClient.token, PostUid: uid},
						)
						if err != nil {
							panic(err)
						}
					case codes.Unavailable:
						time.Sleep(time.Second * 5)
						s.deletePostStatsChannel <- uid
					}
				}
			}
		}()
	}
}

func (s *Server) deleteCommentWorker() {
	ctx := context.Background()
	for {
		uid := <-s.deleteCommentChannel

		log.Printf("DeleteComment rabotyaga: got %s", uid)

		go func() {
			_, err := s.commentClient.client.DeleteComment(ctx,
				&comment.DeleteCommentRequest{Uid: uid, Token: s.commentClient.token},
			)
			if err != nil {
				if st, ok := status.FromError(err); ok {
					switch st.Code() {
					case codes.Unauthenticated:
						err := s.updatePostToken()
						if err != nil {
							panic(err)
						}
						_, err = s.commentClient.client.DeleteComment(ctx,
							&comment.DeleteCommentRequest{Uid: uid, Token: s.commentClient.token},
						)
						if err != nil {
							panic(err)
						}
					case codes.Unavailable:
						time.Sleep(time.Second * 5)
						s.deleteCommentChannel <- uid
					}
				}
			}
		}()
	}
}
