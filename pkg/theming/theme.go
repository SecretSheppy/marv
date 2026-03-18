package theming

import (
	"fmt"
	"strings"
)

type CSS interface {
	ToCSS() string
}

type TextStyle struct {
	Caret, Background, Color, FontStyle, TextDecoration, SelectionBackground, SelectionColor, SelectionBorder string
}

func (t TextStyle) ToCSS() string {
	css := strings.Builder{}

	appendStyle := func(property, value string) {
		if value != "" {
			css.WriteString(property)
			css.WriteString(":")
			css.WriteString(value)
			css.WriteString(";")
		}
	}

	appendStyle("caret-color", t.Caret)
	appendStyle("background", t.Background)
	appendStyle("color", t.Color)
	appendStyle("font-style", t.FontStyle)
	appendStyle("text-decoration", t.TextDecoration)

	if t.SelectionBackground != "" || t.SelectionColor != "" || t.SelectionBorder != "" {
		css.WriteString("&::selection{")
		appendStyle("background", t.SelectionBackground)
		appendStyle("color", t.SelectionColor)
		appendStyle("border", t.SelectionBorder)
		css.WriteString("}")
	}

	return css.String()
}

type Code struct {
	Variable, Function, Keyword, String, Comment, Number, Operator, Type, Class, Interface, Constant, Property TextStyle
}

func (c Code) ToCSS() string {
	css := strings.Builder{}

	appendStyle := func(class, style string) {
		css.WriteString(".")
		css.WriteString(class)
		css.WriteString("{")
		css.WriteString(style)
		css.WriteString("}")
	}

	appendStyle("variable", c.Variable.ToCSS())
	appendStyle("function", c.Function.ToCSS())
	appendStyle("keyword", c.Keyword.ToCSS())
	appendStyle("string", c.String.ToCSS())
	appendStyle("comment", c.Comment.ToCSS())
	appendStyle("number", c.Number.ToCSS())
	appendStyle("operator", c.Operator.ToCSS())
	appendStyle("type", c.Type.ToCSS())
	appendStyle("class", c.Class.ToCSS())
	appendStyle("interface", c.Interface.ToCSS())
	appendStyle("constant", c.Constant.ToCSS())
	appendStyle("property", c.Property.ToCSS())

	return css.String()
}

type InterfaceColors struct {
	// TODO: add when css exists, needs to overwrite :root variables with theme
}

func (i InterfaceColors) ToCSS() string {
	return ":root{}"
}

type Theme struct {
	Colors InterfaceColors
	Code   Code
}

func (t Theme) ToCSS() string {
	return fmt.Sprintf("%s%s", t.Colors.ToCSS(), t.Code.ToCSS())
}
