package prosemirror2html

import (
	"errors"
	"fmt"
)

// Tag is an interface describing how a node or mark is rendered as HTML.
// The Renderer type will pass the node's/mark's attributes to both methods in case you require them,
// for example to render closing heading tags like `</h1>`.
// If either method implementation returns an error, rendering is halted and the error is returned from the
// Render() or RenderNode() call.
type Tag interface {
	// RenderOpening renders the type's opening tag. Return an empty string to render no tags for your custom type.
	RenderOpening(attrs map[string]interface{}) (string, error)
	// RenderClosing renders the type's closing tag. Return an empty string to render no tags for your custom type or if the tag is self-closing.
	RenderClosing(attrs map[string]interface{}) (string, error)
}

// SimpleTag is a simple tag, using Name as tag name and rendering all attributes into the opening tag.
// For example, the default implementation for Prosemirror's 'link' mark type is a SimpleTag, and would
// render 'link' marks as something like `<a href="https://wikipedia.org/" target="_blank">...</a>`.
type SimpleTag struct {
	Name        string
	SelfClosing bool // if true, RenderClosing returns an empty string
}

var _ Tag = SimpleTag{} // compile time 'implements' check

// RenderOpening renders a standard HTML opening tag, with all attributes as `<name>="<value>"`.
// If an attribute value is a number or a boolean, it will omit the surrounding quotes.
func (t SimpleTag) RenderOpening(attrs map[string]interface{}) (string, error) {
	formattedAttrs := ""
	for name, value := range attrs {
		switch v := value.(type) {
		case bool, float64:
			formattedAttrs += fmt.Sprintf(` %s=%v`, name, v)
		default:
			formattedAttrs += fmt.Sprintf(` %s="%v"`, name, v)
		}
	}

	return "<" + t.Name + formattedAttrs + ">", nil
}

// RenderClosing renders a standard HTML closing tag, except if t.SelfClosing is true, in which case
// it renders an empty string.
func (t SimpleTag) RenderClosing(map[string]interface{}) (string, error) {
	if t.SelfClosing {
		return "", nil
	}

	return "</" + t.Name + ">", nil
}

// default implementation for Prosemirror's text nodes.
type text struct{}

var _ Tag = text{} // compile time 'implements' check

func (text) RenderOpening(map[string]interface{}) (string, error) { return "", nil }
func (text) RenderClosing(map[string]interface{}) (string, error) { return "", nil }

type heading struct{}

var _ Tag = heading{} // compile time 'implements' check

func (heading) RenderOpening(attrs map[string]interface{}) (string, error) {
	levelRaw, ok := attrs["level"]
	if !ok {
		return "", errors.New("heading has no level attribute")
	}

	levelFloat, ok := levelRaw.(float64) // float64 because of json.Unmarshal
	if !ok {
		return "", errors.New("heading's level attribute is not a number")
	}
	level := int(levelFloat)

	return fmt.Sprintf("<h%d>", level), nil
}

func (heading) RenderClosing(attrs map[string]interface{}) (string, error) {
	levelRaw, ok := attrs["level"]
	if !ok {
		return "", errors.New("heading has no level attribute")
	}

	levelFloat, ok := levelRaw.(float64) // float64 because of json.Unmarshal
	if !ok {
		return "", errors.New("heading's level attribute is not a number")
	}
	level := int(levelFloat)

	return fmt.Sprintf("</h%d>", level), nil
}
