package server

import (
	"github.com/gorilla/mux"
	"github.com/rdsalakhov/xsolla-login-test/services"
	"net/http"
)

type server struct {
	router      *mux.Router
	authService services.AuthServiceProvider
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) Start(port string) {
	http.ListenAndServe(port, s.router)
}

func (s *server) ConfigureRouter() {
	s.router.HandleFunc("/", s.handleHome)
	s.router.HandleFunc("/callback", s.handleCallback)
}

func NewServer() *server {
	authService := services.NewAuthService()
	server := &server{
		router:      mux.NewRouter(),
		authService: authService,
	}

	server.ConfigureRouter()
	return server
}




