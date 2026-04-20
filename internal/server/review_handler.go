package server

import (
	"encoding/json"
	"net/http"

	"github.com/SecretSheppy/marv/internal/review"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type reviewRequest struct {
	File   string `json:"file"`
	Review string `json:"review"`
}

func (s *Server) reviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fwName := vars["framework"]
	framework := s.getActiveFw(fwName)
	if framework == nil {
		writeAPIError(w, r, nil, http.StatusNotFound, "error: no active framework with name "+fwName)
	}
	mutantID, err := uuid.Parse(vars["mutant-id"])
	if err != nil {
		writeAPIError(w, r, err, http.StatusInternalServerError, "error: failed to parse mutant id")
		return
	}

	var reviewReq reviewRequest
	if err := json.NewDecoder(r.Body).Decode(&reviewReq); err != nil {
		writeAPIError(w, r, err, http.StatusBadRequest, "error: invalid request format, failed to decode request")
		return
	}

	_, mutation := framework.Mutations()[reviewReq.File].GetMutant(mutantID)
	fwid := ""
	if mutation != nil {
		fwid = mutation.FrameworkMutantID
	}
	rev := &review.Review{
		MutationID:          mutantID,
		FrameworkMutationID: fwid,
		Framework:           fwName,
		Review:              reviewReq.Review,
	}
	if err := s.db.SaveReview(rev); err != nil {
		writeAPIError(w, r, err, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}
