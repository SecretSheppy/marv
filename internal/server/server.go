package server

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Server struct {
	port int
	fws  []fwlib.Framework
}

func NewServer(port int, fws []fwlib.Framework) *Server {
	return &Server{port: port, fws: fws}
}

func (s *Server) Serve() error {
	r := mux.NewRouter()
	r.Use(logger)
	r.PathPrefix("/{framework}/mutant/").HandlerFunc(s.mutantHandler).Methods("GET")
	r.PathPrefix("/{framework}/mutants/").HandlerFunc(s.mutantsHandler).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", s.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv.ListenAndServe()
}

func (s *Server) getActiveFw(fwName string) fwlib.Framework {
	for _, fw := range s.fws {
		if fw.Meta().Name == fwName {
			return fw
		}
	}
	return nil
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().Str("method", r.Method).Str("uri", r.RequestURI).Msg("New request")
		next.ServeHTTP(w, r)
	})
}

func writeError(w http.ResponseWriter, r *http.Request, err error, code int, message string) {
	log.Error().Err(err).Str("path", r.URL.Path).Msg(message)
	w.WriteHeader(http.StatusBadRequest)
	var buff bytes.Buffer
	buff.WriteString("<html><body>")
	buff.WriteString(fmt.Sprintf("<h1>%d: %s</h1>", code, http.StatusText(code)))
	buff.WriteString(fmt.Sprintf("<p>%s</p>", message))
	buff.WriteString("</body></html>")
	w.Write(buff.Bytes())
}
