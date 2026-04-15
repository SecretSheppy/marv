package html

import (
	"os"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

type Meta struct {
	StylePaths, ScriptPaths []string
	style                   []byte
	scripts                 [][]byte
}

// MinifyAndCache minifies and caches all styles and scripts. All stylesheets are reduced down into a single entity,
// whereas scripts are kept independent.
func (m *Meta) MinifyAndCache() error {
	mini := minify.New()
	mini.AddFunc("text/css", css.Minify)
	mini.AddFunc("text/javascript", js.Minify)

	if err := m.minifyAndCacheCss(mini); err != nil {
		return err
	}
	return m.minifyAndCacheJs(mini)
}

// compiles all styles into one large minified style
func (m *Meta) minifyAndCacheCss(mini *minify.M) error {
	for _, style := range m.StylePaths {
		content, err := os.ReadFile(style)
		if err != nil {
			return err
		}
		minified, err := mini.Bytes("text/css", content)
		if err != nil {
			return err
		}
		m.style = append(m.style, minified...)
	}
	return nil
}

func (m *Meta) minifyAndCacheJs(mini *minify.M) error {
	for _, script := range m.ScriptPaths {
		content, err := os.ReadFile(script)
		if err != nil {
			return err
		}
		minified, err := mini.Bytes("text/javascript", content)
		if err != nil {
			return err
		}
		m.scripts = append(m.scripts, minified)
	}
	return nil
}

func (m *Meta) Style() []byte {
	return m.style
}

func (m *Meta) Scripts() [][]byte {
	return m.scripts
}
