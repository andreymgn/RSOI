package post

import (
	"errors"
	"testing"
	"time"

	pb "github.com/andreymgn/RSOI/services/post/proto"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

var (
	errDummy = errors.New("dummy")
)

type mockdb struct{}

func (mdb *mockdb) getAll(pageSize, pageNumber int32) ([]*Post, error) {
	result := make([]*Post, 0)
	result = append(result, &Post{"dba6af76-fc50-4aab-a12a-7afd1c1b2330", "First post", "google.com", time.Now(), time.Now()})
	result = append(result, &Post{"4007dd37-ecc2-4f00-a7ba-457b9f3e7eb2", "Second post", "", time.Now(), time.Now().Add(time.Second * 10)})
	result = append(result, &Post{"dba6af76-fc50-4aab-a12a-7afd1c1b2330", "Third post", "yandex.ru", time.Now(), time.Now()})
	return result, nil
}

func (mdb *mockdb) getOne(uid string) (*Post, error) {
	if uid == "success" {
		return &Post{"dba6af76-fc50-4aab-a12a-7afd1c1b2330", "First post", "google.com", time.Now(), time.Now()}, nil
	} else {
		return nil, errDummy
	}
}

func (mdb *mockdb) create(title, url string) error {
	if title == "success" {
		return nil
	} else {
		return errDummy
	}
}

func (mdb *mockdb) update(uid, title, url string) error {
	if uid == "success" {
		return nil
	} else {
		return errDummy
	}
}

func (mdb *mockdb) delete(uid string) error {
	if uid == "success" {
		return nil
	} else {
		return errDummy
	}
}

func TestListPosts(t *testing.T) {
	s := &Server{&mockdb{}}
	var pageSize int32 = 3
	req := &pb.ListPostsRequest{PageSize: pageSize, PageNumber: 1}
	res, err := s.ListPosts(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(res.Posts) != int(pageSize) {
		t.Errorf("unexpected number of posts: got %v want %v", len(res.Posts), pageSize)
	}

	if res.Posts[0].Uid != "dba6af76-fc50-4aab-a12a-7afd1c1b2330" {
		t.Errorf("unexpected post uid")
	}
}

func TestGetPost(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.GetPostRequest{Uid: "success"}
	_, err := s.GetPost(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestGetPostFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.GetPostRequest{Uid: ""}
	_, err := s.GetPost(context.Background(), req)
	if err != ErrUidNotSet {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.GetPostRequest{Uid: "fail"}
	_, err = s.GetPost(context.Background(), req)
	if err != errDummy {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCreatePost(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.CreatePostRequest{Title: "success"}
	_, err := s.CreatePost(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCreatePostFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.CreatePostRequest{Title: ""}
	_, err := s.CreatePost(context.Background(), req)
	if err != ErrTitleNotSet {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.CreatePostRequest{Title: "fail"}
	_, err = s.CreatePost(context.Background(), req)
	if err != errDummy {
		t.Errorf("unexpected error %v", err)
	}
}

func TestUpdatePost(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.UpdatePostRequest{Uid: "success"}
	_, err := s.UpdatePost(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestUpdatePostFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.UpdatePostRequest{Uid: ""}
	_, err := s.UpdatePost(context.Background(), req)
	if err != ErrUidNotSet {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.UpdatePostRequest{Uid: "fail"}
	_, err = s.UpdatePost(context.Background(), req)
	if err != errDummy {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDeletePost(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.DeletePostRequest{Uid: "success"}
	_, err := s.DeletePost(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDeletePostFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.DeletePostRequest{Uid: ""}
	_, err := s.DeletePost(context.Background(), req)
	if err != ErrUidNotSet {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.DeletePostRequest{Uid: "fail"}
	_, err = s.DeletePost(context.Background(), req)
	if err != errDummy {
		t.Errorf("unexpected error %v", err)
	}
}
