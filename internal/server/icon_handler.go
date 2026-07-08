package server

import (
	"io/fs"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/SecretSheppy/marv/web"
	"github.com/gorilla/mux"
)

var re = regexp.MustCompile("fill=\"(#[a-zA-Z0-9]+)\"")

func (s *Server) iconHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	color := vars["color"]
	name := vars["name"]

	if strings.HasPrefix(color, "#") {
		color = color[1:]
	}

	file, err := fs.ReadFile(web.IconsFS, path.Join("icons", name))
	if err != nil {
		writeAPIError(w, r, nil, http.StatusNotFound, "icon "+name+" not found")
		return
	}

	file = re.ReplaceAll(file, []byte("fill=\"#"+color+"\""))

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(file)
}
