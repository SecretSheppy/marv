package themes

import (
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/SecretSheppy/marv/pkg/colour"
)

const (
	Logo     = "/resources/branding/marv_logo.png"
	LogoDark = "/resources/branding/marv_logo_dark.png"
)

type Document struct {
	Background string `json:"background"`
}

func (d Document) CSS() string {
	return fmt.Sprintf("--body-bg:%s;", d.Background)
}

type UITextColor struct {
	Main string `json:"main"`
	Dark string `json:"dark"`
}

func (u UITextColor) CSS() string {
	return fmt.Sprintf("--main-text-fg:%s;--darker-text-fg:%s;", u.Main, u.Dark)
}

type UIText struct {
	FontFamily string      `json:"font-family"`
	TitleSize  string      `json:"title-size"`
	FontSize   string      `json:"font-size"`
	Color      UITextColor `json:"color"`
}

func (u UIText) CSS() string {
	return fmt.Sprintf("--ui-font-family:%s;--ui-title-size:%s;--ui-font-size:%s;%s",
		u.FontFamily, u.TitleSize, u.FontSize, u.Color.CSS())
}

type UIColorsAccent struct {
	HoverBackground string `json:"hover-background"`
	FocusBackground string `json:"focus-background"`
}

func (u UIColorsAccent) CSS() string {
	return fmt.Sprintf("--theme-hover-bg:%s;--theme-focus-bg:%s;", u.HoverBackground, u.FocusBackground)
}

type UIColorsStatistics struct {
	Positive string `json:"positive"`
	Mediocre string `json:"mediocre"`
	Negative string `json:"negative"`
}

func (u UIColorsStatistics) CSS() string {
	return fmt.Sprintf("--stat-green-fg:%s;--stat-orange-fg:%s;--stat-red-fg:%s;",
		u.Positive, u.Mediocre, u.Negative)
}

type UIColors struct {
	PrimaryBackground   string             `json:"primary-background"`
	SecondaryBackground string             `json:"secondary-background"`
	HoverBackground     string             `json:"hover-background"`
	Accent              UIColorsAccent     `json:"accent"`
	Statistics          UIColorsStatistics `json:"statistics"`
	LinkColor           string             `json:"link-color"`
}

func (u UIColors) CSS() string {
	return fmt.Sprintf("--main-bg:%s;--main-secondary-bg:%s;--main-hover-bg:%s;%s%s--link-fg:%s;",
		u.PrimaryBackground, u.SecondaryBackground, u.HoverBackground, u.Accent.CSS(), u.Statistics.CSS(), u.LinkColor)
}

type UIComponents struct {
	Divider      string `json:"divider"`
	HeaderHeight string `json:"header-height"`
}

func (u UIComponents) CSS() string {
	return fmt.Sprintf("--main-divider:%s;--main-header-height:%s;", u.Divider, u.HeaderHeight)
}

type UI struct {
	Text       UIText       `json:"text"`
	Colors     UIColors     `json:"colors"`
	Components UIComponents `json:"components"`
}

func (u UI) CSS() string {
	return fmt.Sprintf("%s%s%s", u.Text.CSS(), u.Colors.CSS(), u.Components.CSS())
}

type CodeText struct {
	FontFamily string `json:"font-family"`
	FontSize   string `json:"font-size"`
}

func (c CodeText) CSS() string {
	return fmt.Sprintf("--code-font-family:%s;--code-font-size:%s;", c.FontFamily, c.FontSize)
}

type CodeDiffColors struct {
	Default string `json:"default"`
	Dark    string `json:"dark"`
	Inline  string `json:"inline"`
}

func (c CodeDiffColors) CSS(prefix string) string {
	return fmt.Sprintf("--diff-%s:%s;--diff-%s-dark:%s;--diff-inline-%s:%s;",
		prefix, c.Default, prefix, c.Dark, prefix, c.Inline)
}

type CodeDiff struct {
	Remove CodeDiffColors `json:"remove"`
	Insert CodeDiffColors `json:"insert"`
}

func (c CodeDiff) CSS() string {
	return fmt.Sprintf("%s%s", c.Remove.CSS("remove"), c.Insert.CSS("insert"))
}

type CodeComponents struct {
	MutationBorder   string `json:"mutation-border"`
	LineNumberBorder string `json:"line-number-border"`
}

func (c CodeComponents) CSS() string {
	return fmt.Sprintf("--mutation-border:%s;--line-number-border:%s;", c.MutationBorder, c.LineNumberBorder)
}

type Code struct {
	ChromaTheme     string         `json:"chroma-theme"`
	Text            CodeText       `json:"text"`
	Background      string         `json:"background"`
	LineNumberColor string         `json:"line-number-color"`
	LineHeight      string         `json:"line-height"`
	Diff            CodeDiff       `json:"diff"`
	Components      CodeComponents `json:"components"`
}

func (c Code) CSS() string {
	return fmt.Sprintf("%s--code-bg:%s;--code-line-number-color:%s;--code-line-height:%s;%s%s",
		c.Text.CSS(), c.Background, c.LineNumberColor, c.LineHeight, c.Diff.CSS(), c.Components.CSS())
}

type StatusColors struct {
	Foreground string `json:"foreground"`
	Background string `json:"background"`
}

func (s StatusColors) CSS(prefix string) string {
	return fmt.Sprintf("--status-%s-fg:%s;--status-%s-bg:%s;", prefix, s.Foreground, prefix, s.Background)
}

type Status struct {
	Killed     StatusColors `json:"killed"`
	Survived   StatusColors `json:"survived"`
	Crashed    StatusColors `json:"crashed"`
	Timeout    StatusColors `json:"timeout"`
	NoCoverage StatusColors `json:"no-coverage"`
}

func (s Status) CSS() string {
	return fmt.Sprintf("%s%s%s%s%s", s.Killed.CSS("killed"), s.Survived.CSS("survived"),
		s.Crashed.CSS("crashed"), s.Timeout.CSS("timeout"), s.NoCoverage.CSS("no-coverage"))
}

type Theme struct {
	Name     string   `json:"name"`
	Document Document `json:"document"`
	UI       UI       `json:"ui"`
	Code     Code     `json:"code"`
	Status   Status   `json:"status"`
	css      string
}

func (t *Theme) CSS() string {
	if t.css == "" {
		t.css = fmt.Sprintf(":root{%s%s%s%s}", t.Document.CSS(), t.UI.CSS(), t.Code.CSS(), t.Status.CSS())
	}
	return t.css
}

func (t *Theme) Logo() string {
	if isBright, _ := colour.IsBright(t.UI.Colors.PrimaryBackground); isBright {
		return LogoDark
	}
	return Logo
}

func (t *Theme) IconColor() string {
	return t.UI.Text.Color.Main[1:]
}

func (t *Theme) Icon(name string) string {
	return fmt.Sprintf("/icon/%s/%s", t.IconColor(), name)
}

func LoadTheme(file string, fsys fs.FS) (*Theme, error) {
	raw, err := fs.ReadFile(fsys, file)
	if err != nil {
		return nil, err
	}
	theme := &Theme{}
	if err := json.Unmarshal(raw, theme); err != nil {
		return nil, err
	}
	return theme, nil
}
