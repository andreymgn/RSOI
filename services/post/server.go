package post

import (
	"errors"
	"fmt"
	"log"
	"net"

	pb "github.com/andreymgn/RSOI/services/post/proto"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	ErrUidNotSet   = errors.New("post UUID is required")
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

// GetPostResponse converts Post to GetPostResponse
func (p *Post) GetPostResponse() (*pb.GetPostResponse, error) {
	createdAtProto, err := ptypes.TimestampProto(p.CreatedAt)
	if err != nil {
		return nil, err
	}

	modifiedAtProto := &pb.NullableTime{}
	if p.ModifiedAt.Valid {
		timeModified, err := ptypes.TimestampProto(p.ModifiedAt.Time)
		if err != nil {
			return nil, err
		}
		modifiedAtProto.Time = timeModified
		modifiedAtProto.Valid = true
	} else {
		modifiedAtProto.Valid = false
	}

	res := new(pb.GetPostResponse)
	res.Uid = p.Uid
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
		postResponse, err := post.GetPostResponse()
		if err != nil {
			return nil, err
		}

		res.Posts = append(res.Posts, postResponse)
	}

	return res, nil
}

// GetPost returns single post by ID
func (s *Server) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	if req.Uid == "" {
		return nil, ErrUidNotSet
	}

	post, err := s.db.getOne(req.Uid)
	if err != nil {
		return nil, err
	}

	return post.GetPostResponse()
}

// CreatePost creates a new post
func (s *Server) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	if req.Title == "" {
		return nil, ErrTitleNotSet
	}
	err := s.db.create(req.Title, req.Url)
	if err != nil {
		return nil, err
	}

	res := new(pb.CreatePostResponse)
	return res, nil
}

// UpdatePost updates post by ID
func (s *Server) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	if req.Uid == "" {
		return nil, ErrUidNotSet
	}

	err := s.db.update(req.Title, req.Url, req.Uid)
	if err != nil {
		return nil, err
	}

	res := new(pb.UpdatePostResponse)
	return res, nil
}

// DeletePost deletes post by ID
func (s *Server) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	if req.Uid == "" {
		return nil, ErrUidNotSet
	}

	err := s.db.delete_(req.Uid)
	if err != nil {
		return nil, err
	}

	res := new(pb.DeletePostResponse)
	return res, nil
}
