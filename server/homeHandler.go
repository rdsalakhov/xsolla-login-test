package server

import (
	"net/http"
)

func (s *server) handleHome(w http.ResponseWriter, r *http.Request) {
	respondHTML(w, r, "loginwidget.html")
}
