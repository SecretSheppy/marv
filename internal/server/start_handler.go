package server

import "net/http"

func (s *Server) startHandler(w http.ResponseWriter, r *http.Request) {
	render, err := s.renderer.RenderStart()
	if err != nil {
		writeError(w, r, err, http.StatusInternalServerError, "")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(render)
}
