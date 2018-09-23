package comment

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/andreymgn/RSOI/services/comment/proto"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	ErrPostUidNotSet    = errors.New("post UUID is required")
	ErrCommentUidNotSet = errors.New("comment UUID is required")
)

// Comment describes comment to a post
type Comment struct {
	Uid        string
	PostUid    string
	Body       string
	ParentUid  string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// Server implements comments service
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
	pb.RegisterCommentServer(server, s)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return server.Serve(lis)
}

// SingleComment converts Comment to SingleComment
func (c *Comment) SingleComment() (*pb.SingleComment, error) {
	createdAtProto, err := ptypes.TimestampProto(c.CreatedAt)
	if err != nil {
		return nil, err
	}

	modifiedAtProto, err := ptypes.TimestampProto(c.ModifiedAt)
	if err != nil {
		return nil, err
	}

	res := new(pb.SingleComment)
	res.Uid = c.Uid
	res.PostUid = c.PostUid
	res.Body = c.Body
	res.ParentUid = c.ParentUid
	res.CreatedAt = createdAtProto
	res.ModifiedAt = modifiedAtProto

	return res, nil
}

// ListComments returns all comments of post
func (s *Server) ListComments(ctx context.Context, req *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	if req.PostUid == "" {
		return nil, ErrPostUidNotSet
	}

	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	query := "SELECT * FROM comments ORDER BY created_at DESC LIMIT $1, $2"
	lastRecord := req.PageNumber * pageSize
	rows, err := s.db.Query(query, lastRecord, pageSize)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := new(pb.ListCommentsResponse)
	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.Uid, &c.PostUid, &c.Body, &c.ParentUid, &c.CreatedAt, &c.ModifiedAt)
		if err != nil {
			return nil, err
		}

		comment, err := c.SingleComment()
		if err != nil {
			return nil, err
		}
		res.Comments = append(res.Comments, comment)
	}

	return res, nil
}

// CreateComment creates a new comment
func (s *Server) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	query := "INSERT INTO comments (post_uid, body, parent_uid) VALUES ($1, $2, $3)"
	_, err := s.db.Query(query, req.PostUid, req.Body, req.ParentUid)
	if err != nil {
		return nil, err
	}

	res := new(pb.CreateCommentResponse)
	return res, nil
}

// UpdateComment updates comment by ID
func (s *Server) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	if req.Uid == "" {
		return nil, ErrCommentUidNotSet
	}

	query := "UPDATE comments SET body=$1 WHERE uid=$2"
	_, err := s.db.Exec(query, req.Body, req.Uid)
	if err != nil {
		return nil, err
	}

	res := new(pb.UpdateCommentResponse)
	return res, nil
}

// DeleteComment deletes post by ID
func (s *Server) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	if req.Uid == "" {
		return nil, ErrCommentUidNotSet
	}

	query := "DELETE FROM POSTS WHERE uid=$1"
	_, err := s.db.Exec(query, req.Uid)
	if err != nil {
		return nil, err
	}

	res := new(pb.DeleteCommentResponse)
	return res, nil
}
