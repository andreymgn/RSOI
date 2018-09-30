package comment

import (
	"fmt"
	"log"
	"net"

	pb "github.com/andreymgn/RSOI/services/comment/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Server implements comments service
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
	res.Uid = c.UID.String()
	res.PostUid = c.PostUID.String()
	res.Body = c.Body
	res.ParentUid = c.ParentUID.String()
	res.CreatedAt = createdAtProto
	res.ModifiedAt = modifiedAtProto

	return res, nil
}

// ListComments returns all comments of post
func (s *Server) ListComments(ctx context.Context, req *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	postUID, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, err
	}

	comments, err := s.db.getAll(postUID, pageSize, req.PageNumber)
	if err != nil {
		return nil, err
	}

	res := new(pb.ListCommentsResponse)
	for _, comment := range comments {
		singleComment, err := comment.SingleComment()
		if err != nil {
			return nil, err
		}
		res.Comments = append(res.Comments, singleComment)
	}

	return res, nil
}

// CreateComment creates a new comment
func (s *Server) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.SingleComment, error) {
	postUID, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, err
	}

	parentUID, err := uuid.Parse(req.ParentUid)

	comment, err := s.db.create(postUID, req.Body, parentUID)
	if err != nil {
		return nil, err
	}

	return comment.SingleComment()
}

// UpdateComment updates comment by ID
func (s *Server) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, err
	}

	err = s.db.update(uid, req.Body)
	if err != nil {
		return nil, err
	}

	res := new(pb.UpdateCommentResponse)
	return res, nil
}

// DeleteComment deletes post by ID
func (s *Server) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, err
	}

	err = s.db.delete(uid)
	if err != nil {
		return nil, err
	}

	res := new(pb.DeleteCommentResponse)
	return res, nil
}
