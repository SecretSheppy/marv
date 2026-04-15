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
	r.PathPrefix("/resources/").Handler(http.StripPrefix("/resources/", http.FileServer(http.Dir("web/static"))))
	r.HandleFunc("/tree", s.treeHandler).Methods("GET")
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
	buff.WriteString("<html><head><style>*{padding:0;margin:0;box-sizing:border-box;}</style></head>" +
		"<body style=\"height:100%;display:flex;justify-content:center;align-items:center;background:#2b2b2b;\"><div>")
	buff.WriteString(fmt.Sprintf("<h1 style=\"color:#a1afb8;\">%d: %s</h1>", code, http.StatusText(code)))
	buff.WriteString(fmt.Sprintf("<p style=\"color:#a1afb8;\">%s</p>", message))
	buff.WriteString("</div></body></html>")
	w.Write(buff.Bytes())
}
