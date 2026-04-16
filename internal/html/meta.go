package html

import (
	"bytes"
	"fmt"
)

// metaRenderer produces the HTML that should be included in the <head> element of the rendered document.
type metaRenderer struct {
	Title, Favicon string
}

func (m *metaRenderer) Render(buff *bytes.Buffer) {
	m.title(buff)
	m.favicon(buff)
}

func (m *metaRenderer) title(buff *bytes.Buffer) {
	buff.WriteString(fmt.Sprintf("<title>%s</title>", m.Title))
}

func (m *metaRenderer) favicon(buff *bytes.Buffer) {
	buff.WriteString(fmt.Sprintf("<link rel=\"icon\" type=\"image/x-icon\" href=\"%s\">", m.Favicon))
}
