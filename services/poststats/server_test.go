package poststats

import (
	"context"
	"errors"
	"testing"

	pb "github.com/andreymgn/RSOI/services/poststats/proto"
)

var (
	errDummy = errors.New("dummy")
)

type mockdb struct{}

func (mdb *mockdb) get(uid string) (*PostStats, error) {
	if uid == "success" {
		return new(PostStats), nil
	}

	return nil, errDummy
}

func (mdb *mockdb) create(uid string) error {
	if uid == "success" {
		return nil
	}

	return errDummy
}

func (mdb *mockdb) like(uid string) error {
	if uid == "success" {
		return nil
	}

	return errDummy
}

func (mdb *mockdb) dislike(uid string) error {
	if uid == "success" {
		return nil
	}

	return errDummy
}

func (mdb *mockdb) view(uid string) error {
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

func TestGet(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.GetPostStatsRequest{PostUid: "success"}
	_, err := s.GetPostStats(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestGetFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.GetPostStatsRequest{}
	_, err := s.GetPostStats(context.Background(), req)
	if err != ErrPostUidNotSet {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.GetPostStatsRequest{PostUid: "fail"}
	_, err = s.GetPostStats(context.Background(), req)
	if err != errDummy {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCreate(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.CreatePostStatsRequest{PostUid: "success"}
	_, err := s.CreatePostStats(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCreateFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.CreatePostStatsRequest{}
	_, err := s.CreatePostStats(context.Background(), req)
	if err != ErrPostUidNotSet {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.CreatePostStatsRequest{PostUid: "fail"}
	_, err = s.CreatePostStats(context.Background(), req)
	if err != errDummy {
		t.Errorf("unexpected error %v", err)
	}
}

func TestLike(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.LikePostRequest{PostUid: "success"}
	_, err := s.LikePost(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestLikeFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.LikePostRequest{}
	_, err := s.LikePost(context.Background(), req)
	if err != ErrPostUidNotSet {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.LikePostRequest{PostUid: "fail"}
	_, err = s.LikePost(context.Background(), req)
	if err != errDummy {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDislike(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.DislikePostRequest{PostUid: "success"}
	_, err := s.DislikePost(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDislikeFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.DislikePostRequest{}
	_, err := s.DislikePost(context.Background(), req)
	if err != ErrPostUidNotSet {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.DislikePostRequest{PostUid: "fail"}
	_, err = s.DislikePost(context.Background(), req)
	if err != errDummy {
		t.Errorf("unexpected error %v", err)
	}
}

func TestView(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.IncreaseViewsRequest{PostUid: "success"}
	_, err := s.IncreaseViews(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestViewFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.IncreaseViewsRequest{}
	_, err := s.IncreaseViews(context.Background(), req)
	if err != ErrPostUidNotSet {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.IncreaseViewsRequest{PostUid: "fail"}
	_, err = s.IncreaseViews(context.Background(), req)
	if err != errDummy {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDelete(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.DeletePostStatsRequest{PostUid: "success"}
	_, err := s.DeletePostStats(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDeleteFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.DeletePostStatsRequest{}
	_, err := s.DeletePostStats(context.Background(), req)
	if err != ErrPostUidNotSet {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.DeletePostStatsRequest{PostUid: "fail"}
	_, err = s.DeletePostStats(context.Background(), req)
	if err != errDummy {
		t.Errorf("unexpected error %v", err)
	}
}
