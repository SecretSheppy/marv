package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) mutantHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fwName := vars["framework"]
	fwPrefix := fmt.Sprintf("/%s/mutant/", fwName)
	file := r.URL.Path[len(fwPrefix):]
	framework := s.getActiveFw(fwName)
	if framework == nil {
		writeError(w, r, nil, http.StatusNotFound, fmt.Sprintf("No active framework with name: %s", fwName))
		return
	}

	mutantID, err := uuid.Parse(r.URL.Query().Get("m"))
	if err != nil {
		writeError(w, r, err, http.StatusBadRequest, "Malformed mutant ID")
		return
	}

	render, err := s.renderer.RenderMutant(framework, file, mutantID)
	switch true {
	case errors.Is(err, ErrFailedToInitRender):
		writeError(w, r, err, http.StatusInternalServerError, err.Error())
		return
	case err != nil:
		writeError(w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(render)
}

func (s *Server) mutantsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fwName := vars["framework"]
	fwPrefix := fmt.Sprintf("/%s/mutants/", fwName)
	file := r.URL.Path[len(fwPrefix):]

	framework := s.getActiveFw(fwName)
	if framework == nil {
		writeError(w, r, nil, http.StatusNotFound, fmt.Sprintf("No active framework with name: %s", fwName))
		return
	}

	render, err := s.renderer.RenderMutants(framework, file)
	switch true {
	case errors.Is(err, ErrFailedToInitRender):
		writeError(w, r, err, http.StatusInternalServerError, err.Error())
		return
	case err != nil:
		writeError(w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(render)
}
