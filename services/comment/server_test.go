package comment

import (
	"errors"
	"testing"
	"time"

	pb "github.com/andreymgn/RSOI/services/comment/proto"
	"golang.org/x/net/context"
)

var (
	errDummy = errors.New("dummy")
)

type mockdb struct{}

func (mdb *mockdb) getAll(postUid string, pageNumber, pageSize int32) ([]*Comment, error) {
	result := make([]*Comment, 0)
	result = append(result, &Comment{"comment-uid-1", "post-uid-1", "first comment body", "", time.Now(), time.Now()})
	result = append(result, &Comment{"comment-uid-2", "post-uid-1", "second comment body", "", time.Now(), time.Now()})
	result = append(result, &Comment{"comment-uid-3", "post-uid-1", "third comment body", "comment-uid-1", time.Now(), time.Now()})
	return result, nil
}

func (mdb *mockdb) create(postUid, body, parentUid string) error {
	if postUid == "success" {
		return nil
	}

	return errDummy
}

func (mdb *mockdb) update(uid, body string) error {
	if uid == "success" {
		return nil
	}

	return errDummy
}

func (mdb *mockdb) delete(uid string) error {
	if uid == "success" {
		return nil
	}

	return errDummy
}

func TestListComments(t *testing.T) {
	s := &Server{&mockdb{}}
	var pageSize int32 = 3
	req := &pb.ListCommentsRequest{PostUid: "post-uid-1", PageSize: pageSize}
	res, err := s.ListComments(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(res.Comments) != int(pageSize) {
		t.Errorf("unexpected number of comments: got %v want %v", len(res.Comments), pageSize)
	}
}

func TestCreateComment(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.CreateCommentRequest{PostUid: "success"}
	_, err := s.CreateComment(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCreateCommentFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.CreateCommentRequest{}
	_, err := s.CreateComment(context.Background(), req)
	if err != ErrPostUidNotSet {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.CreateCommentRequest{PostUid: "fail"}
	_, err = s.CreateComment(context.Background(), req)
	if err != errDummy {
		t.Errorf("unexpected error %v", err)
	}
}

func TestUpdateComment(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.UpdateCommentRequest{Uid: "success"}
	_, err := s.UpdateComment(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestUpdateCommentFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.UpdateCommentRequest{}
	_, err := s.UpdateComment(context.Background(), req)
	if err != ErrCommentUidNotSet {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.UpdateCommentRequest{Uid: "fail"}
	_, err = s.UpdateComment(context.Background(), req)
	if err != errDummy {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDeleteComment(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.DeleteCommentRequest{Uid: "success"}
	_, err := s.DeleteComment(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDeleteCOmmentFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.DeleteCommentRequest{}
	_, err := s.DeleteComment(context.Background(), req)
	if err != ErrCommentUidNotSet {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.DeleteCommentRequest{Uid: "fail"}
	_, err = s.DeleteComment(context.Background(), req)
	if err != errDummy {
		t.Errorf("unexpected error %v", err)
	}
}
