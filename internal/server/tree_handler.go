package server

import (
	"bytes"
	"net/http"

	"github.com/SecretSheppy/marv/internal/html"
)

var treePageMeta = &html.Meta{
	StylePaths:  []string{"web/styles/main.css", "web/styles/tree.css"},
	ScriptPaths: []string{"web/scripts/tree.js"},
}

func (s *Server) treeHandler(w http.ResponseWriter, r *http.Request) {
	if err := treePageMeta.MinifyAndCache(); err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	var buff bytes.Buffer
	buff.WriteString("<html><head><link rel=\"icon\" type=\"image/x-icon\" href=\"/resources/branding/marv_favicon.png\"><style>")
	buff.Write(treePageMeta.Style())
	buff.WriteString("</style>")
	for _, script := range treePageMeta.Scripts() {
		buff.WriteString("<script type=\"text/javascript\">")
		buff.Write(script)
		buff.WriteString("</script>")
	}
	buff.WriteString("</head><body>")
	html.NewTreeRenderer(s.fws).Render(&buff)
	buff.WriteString("</body></html>")
	w.Write(buff.Bytes())
}
