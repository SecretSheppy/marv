package server

import (
	"bytes"
	"net/http"

	"github.com/SecretSheppy/marv/internal/html"
)

var treePageMeta = &html.Meta{
	StylePaths: []string{"web/styles/main.css", "web/styles/tree.css"},
}

func (s *Server) treeHandler(w http.ResponseWriter, r *http.Request) {
	if err := treePageMeta.MinifyAndCache(); err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	var buff bytes.Buffer
	buff.WriteString("<html><head><style>")
	buff.Write(treePageMeta.Style())
	buff.WriteString("</style></head><body>")
	html.NewTreeRenderer(s.fws).Render(&buff)
	buff.WriteString("</body></html>")
	w.Write(buff.Bytes())
}
