package server

import (
	"net/http"
)

func (s *Server) treeHandler(w http.ResponseWriter, r *http.Request) {
	render, err := s.renderer.RenderTree()
	if err != nil {
		writeError(w, r, err, http.StatusInternalServerError, "")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(render)
}
