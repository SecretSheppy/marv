package html

import (
	"bytes"
	"fmt"

	"github.com/SecretSheppy/marv/internal/mutations"
)

func writeFilters(buff *bytes.Buffer) {
	buff.WriteString("<div class=\"filters-component collapsed\">" +
		"<div id=\"filters-toggle\" class=\"filters-bar\">" +
		"<img class=\"icon\" src=\"/resources/icons/sliders-solid.svg\" alt=\"filters icon\" />" +
		"<h4 class=\"bar-title\">Status Filters</h4>" +
		"<div class=\"right-content\">" +
		"<img class=\"icon arrow-up\" src=\"/resources/icons/arrow_up.png\" alt=\"arrow up icon\" />" +
		"<img class=\"icon arrow-down\" src=\"/resources/icons/arrow_down.png\" alt=\"arrow down icon\" />" +
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
			"<input id=\"show-%s\" type=\"checkbox\" checked /> %s"+
			"</label>", status.Text(), status.Text(), status.Text()))
	}
	buff.WriteString("</div>" + // closes filters-wrapper
		"</div>" + // closes content-wrapper
		"</div>")
}
