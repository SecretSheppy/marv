package html

import (
	"bytes"
	"fmt"

	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/SecretSheppy/marv/internal/themes"
)

func writeFilters(buff *bytes.Buffer, theme *themes.Theme) {
	buff.WriteString("<div id=\"filters\" class=\"filters-component collapsed\">" +
		"<div id=\"filters-toggle\" class=\"filters-bar\">" +
		"<img class=\"icon\" src=\"" + getIconURL(theme, "sliders-solid.svg") + "\" alt=\"filters icon\" />" +
		"<h4 class=\"bar-title\">Status Filters</h4>" +
		"<div class=\"right-content\">" +
		"<img class=\"icon arrow-up\" src=\"" + getIconURL(theme, "arrow-up.svg") + "\" alt=\"arrow up icon\" />" +
		"<img class=\"icon arrow-down\" src=\"" + getIconURL(theme, "arrow-down.svg") + "\" alt=\"arrow down icon\" />" +
		"</div>" + // closes right-content
		"</div>" + // closes filters-bar
		"<div class=\"content-wrapper\">" +
		"<p class=\"section-description\">" +
		"Using the filters will hide all mutants with statuses that are not enabled. " +
		"This setting syncs across all open tabs." +
		"</p>" +
		"<div class=\"filters-wrapper\">")
	for _, status := range mutations.Statuses {
		buff.WriteString(fmt.Sprintf("<label for=\"show-%s\" class=\"filter\">"+
			"<input id=\"show-%s\" type=\"checkbox\" name=\"%s\" checked /> %s"+
			"</label>", status.Text(), status.Text(), status.Text(), status.Text()))
	}
	buff.WriteString("</div>" + // closes filters-wrapper
		"</div>" + // closes content-wrapper
		"</div>")
}
