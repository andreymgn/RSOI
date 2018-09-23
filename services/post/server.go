package post

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/andreymgn/RSOI/services/post/proto"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	ErrUidNotSet   = errors.New("UUID is required")
	ErrTitleNotSet = errors.New("post title is required")
)

// Post describes a post
type Post struct {
	Uid        string
	Title      string
	URL        string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// Server implements posts service
type Server struct {
	db *sql.DB
}

// NewServer returns a new server
func NewServer(connString string) (*Server, error) {
	db, err := sql.Open("postgres", connString)
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

	modifiedAtProto, err := ptypes.TimestampProto(p.ModifiedAt)
	if err != nil {
		return nil, err
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

	query := "SELECT * FROM posts ORDER BY created_at DESC LIMIT $1, $2"
	lastRecord := req.PageNumber * pageSize
	rows, err := s.db.Query(query, lastRecord, pageSize)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := new(pb.ListPostsResponse)
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.Uid, &p.Title, &p.URL, &p.CreatedAt, &p.ModifiedAt)
		if err != nil {
			return nil, err
		}

		post, err := p.GetPostResponse()
		if err != nil {
			return nil, err
		}

		res.Posts = append(res.Posts, post)
	}

	return res, nil
}

// GetPost returns single post by ID
func (s *Server) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	if req.Uid == "" {
		return nil, ErrUidNotSet
	}

	query := "SELECT * FROM posts WHERE uid=$1"
	row := s.db.QueryRow(query, req.Uid)
	var p Post
	switch err := row.Scan(&p.Uid, &p.Title, &p.URL, &p.CreatedAt, &p.ModifiedAt); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return p.GetPostResponse()
	default:
		return nil, err
	}
}

// CretePost creates a new post
func (s *Server) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	query := "INSERT INTO POSTS (title, url) VALUES ($1, $2)"
	_, err := s.db.Exec(query, req.Title, req.Url)
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

	query := "UPDATE posts SET title=COALESCE(NULLIF($1,''), title), url=COALESCE(NULLIF($2,''), url) WHERE uid=$3"
	_, err := s.db.Exec(query, req.Title, req.Url, req.Uid)
	if err != nil {
		return nil, err
	}

	res := new(pb.UpdatePostResponse)
	return res, nil
}

// DeleltePost deletes post by ID
func (s *Server) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	if req.Uid == "" {
		return nil, ErrUidNotSet
	}

	query := "DELETE FROM posts WHERE uid=$1"
	_, err := s.db.Exec(query, req.Uid)
	if err != nil {
		return nil, err
	}

	res := new(pb.DeletePostResponse)
	return res, nil
}
