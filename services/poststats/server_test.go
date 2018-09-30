package poststats

import (
	"context"
	"errors"
	"testing"

	pb "github.com/andreymgn/RSOI/services/poststats/proto"
	"github.com/google/uuid"
)

var (
	errDummy     = errors.New("dummy")
	dummyUID     = uuid.New()
	nilUIDString = uuid.Nil.String()
)

type mockdb struct{}

func (mdb *mockdb) get(uid uuid.UUID) (*PostStats, error) {
	if uid == uuid.Nil {
		return new(PostStats), nil
	}

	return nil, errDummy
}

func (mdb *mockdb) create(uid uuid.UUID) (*PostStats, error) {
	if uid == uuid.Nil {
		return new(PostStats), nil
	}

	return nil, errDummy
}

func (mdb *mockdb) like(uid uuid.UUID) error {
	if uid == uuid.Nil {
		return nil
	}

	return errDummy
}

func (mdb *mockdb) dislike(uid uuid.UUID) error {
	if uid == uuid.Nil {
		return nil
	}

	return errDummy
}

func (mdb *mockdb) view(uid uuid.UUID) error {
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

func TestGet(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.GetPostStatsRequest{PostUid: nilUIDString}
	_, err := s.GetPostStats(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestGetFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.GetPostStatsRequest{}
	_, err := s.GetPostStats(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestCreate(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.CreatePostStatsRequest{PostUid: nilUIDString}
	_, err := s.CreatePostStats(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCreateFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.CreatePostStatsRequest{}
	_, err := s.CreatePostStats(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestLike(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.LikePostRequest{PostUid: nilUIDString}
	_, err := s.LikePost(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestLikeFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.LikePostRequest{}
	_, err := s.LikePost(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestDislike(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.DislikePostRequest{PostUid: nilUIDString}
	_, err := s.DislikePost(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDislikeFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.DislikePostRequest{}
	_, err := s.DislikePost(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestView(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.IncreaseViewsRequest{PostUid: nilUIDString}
	_, err := s.IncreaseViews(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestViewFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.IncreaseViewsRequest{}
	_, err := s.IncreaseViews(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestDelete(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.DeletePostStatsRequest{PostUid: nilUIDString}
	_, err := s.DeletePostStats(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDeleteFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.DeletePostStatsRequest{}
	_, err := s.DeletePostStats(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}
