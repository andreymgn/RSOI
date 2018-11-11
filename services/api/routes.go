package api

import "net/http"

func (s *Server) routes() {
	postsRouter := s.router.Mux.PathPrefix("/api/posts").Subrouter()
	postsRouter.HandleFunc("/", s.getPosts()).Methods("GET")
	postsRouter.HandleFunc("/", s.createPost()).Methods("POST")
	postsRouter.HandleFunc("/{uid}", s.getPost()).Methods("GET")
	postsRouter.HandleFunc("/{uid}", s.updatePost()).Methods("PATCH")
	postsRouter.HandleFunc("/{uid}", s.deletePost()).Methods("DELETE")

	postsRouter.HandleFunc("/{uid}/like", s.likePost()).Methods("PATCH")
	postsRouter.HandleFunc("/{uid}/dislike", s.dislikePost()).Methods("PATCH")

	postsRouter.HandleFunc("/{postuid}/comments/", s.getPostComments()).Methods("GET")
	postsRouter.HandleFunc("/{postuid}/comments/", s.createComment()).Methods("POST")
	postsRouter.HandleFunc("/{postuid}/comments/{uid}", s.getPostComments()).Methods("GET")
	postsRouter.HandleFunc("/{postuid}/comments/{uid}", s.updateComment()).Methods("PATCH")
	postsRouter.HandleFunc("/{postuid}/comments/{uid}", s.deleteComment()).Methods("DELETE")

	s.router.Mux.HandleFunc("/api/user", s.createUser()).Methods("POST")
	s.router.Mux.HandleFunc("/api/auth", s.getUserToken()).Methods("POST")

	s.router.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Hello, world!")) })
}
