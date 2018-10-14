package comment

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/andreymgn/RSOI/services/comment/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	statusInvalidUUID = status.Error(codes.InvalidArgument, "invalid UUID")
	statusNotFound    = status.Error(codes.NotFound, "comment not found")
)

func internalError(err error) error {
	return status.Error(codes.Internal, err.Error())
}

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
func (s *Server) Start(port int, tracer opentracing.Tracer) error {
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
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
		return nil, internalError(err)
	}

	modifiedAtProto, err := ptypes.TimestampProto(c.ModifiedAt)
	if err != nil {
		return nil, internalError(err)
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
		return nil, statusInvalidUUID
	}

	comments, err := s.db.getAll(postUID, pageSize, req.PageNumber)
	if err != nil {
		return nil, internalError(err)
	}

	res := new(pb.ListCommentsResponse)
	for _, comment := range comments {
		singleComment, err := comment.SingleComment()
		if err != nil {
			return nil, err
		}
		res.Comments = append(res.Comments, singleComment)
	}

	res.PageSize = pageSize
	res.PageNumber = req.PageNumber

	return res, nil
}

// CreateComment creates a new comment
func (s *Server) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.SingleComment, error) {
	postUID, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	parentUID := uuid.Nil
	if req.ParentUid != "" {
		parentUID, err = uuid.Parse(req.ParentUid)
		if err != nil {
			return nil, statusInvalidUUID
		}
	}

	comment, err := s.db.create(postUID, req.Body, parentUID)
	if err != nil {
		return nil, internalError(err)
	}

	return comment.SingleComment()
}

// UpdateComment updates comment by ID
func (s *Server) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.update(uid, req.Body)
	switch err {
	case nil:
		return new(pb.UpdateCommentResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// DeleteComment deletes post by ID
func (s *Server) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, err
	}

	err = s.db.delete(uid)
	switch err {
	case nil:
		return new(pb.DeleteCommentResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}
