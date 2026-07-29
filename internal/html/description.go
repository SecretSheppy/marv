package html

import (
	"bytes"
	"strings"
)

const backtick = 96

func formatDescription(description string) string {
	d := []byte(description)
	var (
		head, tail int
		open       bool
		desc       strings.Builder
		code       bytes.Buffer
	)
	for i := 0; i < len(d); i++ {
		b := d[i]
		switch {
		case b == backtick && head == 0:
			for j := i; j < len(d) && d[j] == backtick; j++ {
				head++
			}
			i += head - 1
			open = true
		case b == backtick:
			for j := i; j < len(d) && d[j] == backtick; j++ {
				tail++
			}
			i += tail - 1
			if head == tail {
				desc.WriteString(`<span class="inline-code">`)
				desc.Write(code.Bytes())
				desc.WriteString(`</span>`)
				code.Reset()
				head = 0
				tail = 0
				open = false
				continue
			}
			tail = 0
			fallthrough
		case open:
			code.Write(d[i : i+1])
		default:
			desc.Write(d[i : i+1])
		}
	}
	return desc.String()
}
