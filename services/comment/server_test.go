package comment

import (
	"errors"
	"testing"
	"time"

	pb "github.com/andreymgn/RSOI/services/comment/proto"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

var (
	errDummy     = errors.New("dummy")
	dummyUID     = uuid.New()
	nilUIDString = uuid.Nil.String()
)

type mockdb struct{}

func (mdb *mockdb) getAll(postUID uuid.UUID, parentUID uuid.UUID, pageNumber, pageSize int32) ([]*Comment, error) {
	result := make([]*Comment, 0)
	uid1 := uuid.New()
	uid2 := uuid.New()
	uid3 := uuid.New()
	pUID := uuid.New()

	result = append(result, &Comment{uid1, pUID, "first comment body", uuid.Nil, time.Now(), time.Now()})
	result = append(result, &Comment{uid2, pUID, "second comment body", uuid.Nil, time.Now(), time.Now()})
	result = append(result, &Comment{uid3, pUID, "third comment body", uid1, time.Now(), time.Now()})
	return result, nil
}

func (mdb *mockdb) create(postUID uuid.UUID, body string, parentUid uuid.UUID) (*Comment, error) {
	if postUID == uuid.Nil {
		uid := uuid.New()
		return &Comment{uid, uid, "first comment body", uuid.Nil, time.Now(), time.Now()}, nil
	}

	return nil, errDummy
}

func (mdb *mockdb) update(uid uuid.UUID, body string) error {
	if uid == uuid.Nil {
		return nil
	}

	return errDummy
}

func (mdb *mockdb) delete(uid uuid.UUID) error {
	if uid == uuid.Nil {
		return nil
	}

	return errDummy
}

func TestListComments(t *testing.T) {
	s := &Server{&mockdb{}}
	var pageSize int32 = 3
	req := &pb.ListCommentsRequest{PostUid: nilUIDString, PageSize: pageSize}
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
	req := &pb.CreateCommentRequest{PostUid: nilUIDString}
	_, err := s.CreateComment(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCreateCommentFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.CreateCommentRequest{}
	_, err := s.CreateComment(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestUpdateComment(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.UpdateCommentRequest{Uid: nilUIDString}
	_, err := s.UpdateComment(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestUpdateCommentFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.UpdateCommentRequest{}
	_, err := s.UpdateComment(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestDeleteComment(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.DeleteCommentRequest{Uid: nilUIDString}
	_, err := s.DeleteComment(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDeleteCOmmentFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.DeleteCommentRequest{}
	_, err := s.DeleteComment(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}
