package server

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/html"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var mutantPageMeta = &html.Meta{
	StylePaths: []string{"web/styles/main.css", "web/styles/code.css"},
}

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
	conflict, mutant := framework.Mutations()[file].GetMutant(mutantID)
	if mutant == nil {
		writeError(w, r, nil, http.StatusNotFound, fmt.Sprintf("No mutant with ID: %s", mutantID))
		return
	}

	conflicts := mutations.Conflicts{
		&mutations.Conflict{
			ID:        conflict.ID,
			StartLine: conflict.StartLine,
			EndLine:   conflict.EndLine,
			Mutations: []*mutations.Mutation{mutant},
		},
	}

	code, err := formatCode(framework, file, conflicts)
	switch true {
	case errors.Is(err, ErrFailedToInitRender):
		writeError(w, r, err, http.StatusInternalServerError, err.Error())
		return
	case err != nil:
		writeError(w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(code)
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
	conflicts := framework.Mutations()[file]

	code, err := formatCode(framework, file, conflicts)
	switch true {
	case errors.Is(err, ErrFailedToInitRender):
		writeError(w, r, err, http.StatusInternalServerError, err.Error())
		return
	case err != nil:
		writeError(w, r, err, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(code)
}

// reads the sources lines of the target file and returns them.
func readLines(srcDir, file string) ([]string, error) {
	data, err := os.ReadFile(path.Join(srcDir, file))
	if err != nil {
		return nil, err
	}
	ls := strings.Lines(string(data))
	lines := make([]string, 0)
	for line := range ls {
		lines = append(lines, strings.ReplaceAll(line, "\n", ""))
	}
	return lines, nil
}

func formatCode(framework fwlib.Framework, file string, conflicts mutations.Conflicts) ([]byte, error) {
	lines, err := readLines(framework.Yaml().SourceCodeDir(), file)
	if err != nil {
		return nil, ErrFailedToReadFile
	}
	if err := mutantPageMeta.MinifyAndCache(); err != nil {
		return nil, ErrFailedToMinifyOrCache
	}
	renderer, err := html.NewRenderer(mutantPageMeta, framework.Meta(), file, lines, conflicts)
	if err != nil {
		return nil, ErrFailedToInitRender
	}
	var buff bytes.Buffer
	if err := renderer.Render(&buff); err != nil {
		return nil, ErrFailedToRenderHTML
	}
	return buff.Bytes(), nil
}
