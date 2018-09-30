package post

import (
	"errors"
	"fmt"
	"log"
	"net"

	pb "github.com/andreymgn/RSOI/services/post/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	ErrTitleNotSet = errors.New("post title is required")
)

// Server implements posts service
type Server struct {
	db datastore
}

// NewServer returns a new server
func NewServer(connString string) (*Server, error) {
	db, err := newDB(connString)
	return &Server{db}, err
}

// Start starts a server
func (s *Server) Start(port int) error {
	server := grpc.NewServer()
	pb.RegisterPostServer(server, s)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return server.Serve(lis)
}

// SinglePost converts Post to SinglePost
func (p *Post) SinglePost() (*pb.SinglePost, error) {
	createdAtProto, err := ptypes.TimestampProto(p.CreatedAt)
	if err != nil {
		return nil, err
	}

	modifiedAtProto, err := ptypes.TimestampProto(p.CreatedAt)
	if err != nil {
		return nil, err
	}

	res := new(pb.SinglePost)
	res.Uid = p.UID.String()
	res.Title = p.Title
	res.Url = p.URL
	res.CreatedAt = createdAtProto
	res.ModifiedAt = modifiedAtProto

	return res, nil
}

// ListPosts returns newest posts
func (s *Server) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	posts, err := s.db.getAll(pageSize, req.PageNumber)
	if err != nil {
		return nil, err
	}
	res := new(pb.ListPostsResponse)
	for _, post := range posts {
		postResponse, err := post.SinglePost()
		if err != nil {
			return nil, err
		}

		res.Posts = append(res.Posts, postResponse)
	}

	return res, nil
}

// GetPost returns single post by ID
func (s *Server) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.SinglePost, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, err
	}

	post, err := s.db.getOne(uid)
	if err != nil {
		return nil, err
	}

	return post.SinglePost()
}

// CreatePost creates a new post
func (s *Server) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.SinglePost, error) {
	if req.Title == "" {
		return nil, ErrTitleNotSet
	}

	post, err := s.db.create(req.Title, req.Url)
	if err != nil {
		return nil, err
	}

	return post.SinglePost()
}

// UpdatePost updates post by ID
func (s *Server) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, err
	}

	err = s.db.update(uid, req.Title, req.Url)
	if err != nil {
		return nil, err
	}

	res := new(pb.UpdatePostResponse)
	return res, nil
}

// DeletePost deletes post by ID
func (s *Server) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, err
	}

	err = s.db.delete(uid)
	if err != nil {
		return nil, err
	}

	res := new(pb.DeletePostResponse)
	return res, nil
}
