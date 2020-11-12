package prosemirror2html

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"strings"
)

// Renderer holds the registered node and mark types.
type Renderer struct {
	nodes map[string][]Tag
	marks map[string][]Tag
}

// NewRenderer returns a new Renderer, with the default node and mark types already registered.
func NewRenderer() *Renderer {
	return &Renderer{
		nodes: map[string][]Tag{
			"text":         {text{}},
			"paragraph":    {SimpleTag{Name: "p"}},
			"blockquote":   {SimpleTag{Name: "blockquote"}},
			"bullet_list":  {SimpleTag{Name: "ul"}},
			"heading":      {heading{}},
			"hard_break":   {SimpleTag{Name: "br", SelfClosing: true}},
			"image":        {SimpleTag{Name: "img", SelfClosing: true}},
			"list_item":    {SimpleTag{Name: "li"}},
			"ordered_list": {SimpleTag{Name: "ol"}},
			"table":        {SimpleTag{Name: "table"}, SimpleTag{Name: "tbody"}},
			"table_cell":   {SimpleTag{Name: "td"}},
			"table_header": {SimpleTag{Name: "th"}},
			"table_row":    {SimpleTag{Name: "tr"}},
		},
		marks: map[string][]Tag{
			"link":        {SimpleTag{Name: "a"}},
			"bold":        {SimpleTag{Name: "strong"}},
			"code":        {SimpleTag{Name: "code"}},
			"italic":      {SimpleTag{Name: "em"}},
			"strike":      {SimpleTag{Name: "s"}},
			"subscript":   {SimpleTag{Name: "sub"}},
			"superscript": {SimpleTag{Name: "sup"}},
			"underline":   {SimpleTag{Name: "u"}},
		},
	}
}

// RegisterNode registers a custom node implementation.
// Registering a custom Tag implementation under the same name as a default node type will
// override the default implementation.
// A custom node type may want to render multiple opening and closing tags
// (like for example the 'table' type, which renders as '<table><tbody>...</tbody></table>'),
// hence why you can pass multiple Tag implementations.
func (r *Renderer) RegisterNode(typ string, tags ...Tag) { r.nodes[typ] = tags }

// RegisterMark registers a custom mark implementation.
// Registering a custom Tag implementation under the same name as a default mark type will
// override the default implementation.
// A custom mark type may want to render multiple opening and closing tags
// hence why you can pass multiple Tag implementations.
func (r *Renderer) RegisterMark(typ string, tags ...Tag) { r.marks[typ] = tags }

// Render parses a Prosemirror JSON document and renders the
// contents using the nodes and mars registered with r.
func (r *Renderer) Render(doc []byte) (string, error) {
	root, err := r.ParseNode(doc)
	if err != nil {
		return "", err
	}

	if root.Type != "doc" {
		return "", errors.New("not a document root node")
	}

	html := []string{}

	for _, n := range root.Content {
		rendered, err := r.RenderNode(n)
		if err != nil {
			return "", err
		}
		html = append(html, rendered)
	}

	return strings.Join(html, ""), nil
}

// ParseNode parses the given Prosemirror JSON to a Node that can be rendered.
func (r *Renderer) ParseNode(data []byte) (*Node, error) {
	n := Node{}

	err := json.Unmarshal(data, &n)
	if err != nil {
		return nil, fmt.Errorf("prosemirror2html: %w", err)
	}

	return &n, nil
}

// RenderNode returns the given node as HTML string.
// RenderNode returns an error when an unknown mark or node is encountered.
// RenderNode renders a nodes content children nodes, if any is given. When no
// content children nodes are found, it renders the nodes text property.
func (r *Renderer) RenderNode(n *Node) (string, error) {
	html := []string{}

	// render opening tags of surrounding marks
	for _, m := range n.Marks {
		tags := r.marks[m.Type]
		if tags == nil {
			return "", fmt.Errorf("unknown mark '%s'", m.Type)
		}
		for _, t := range tags {
			openTag, err := t.RenderOpening(m.Attrs)
			if err != nil {
				return "", fmt.Errorf("prosemirror2html: %w", err)
			}
			html = append(html, openTag)
		}
	}

	// render opening tag(s) of node
	tags, ok := r.nodes[n.Type]
	if !ok {
		return "", fmt.Errorf("unknown node '%s'", n.Type)
	}
	for _, t := range tags {
		openTag, err := t.RenderOpening(n.Attrs)
		if err != nil {
			return "", fmt.Errorf("prosemirror2html: %w", err)
		}
		html = append(html, openTag)
	}

	// render children nodes OR text
	if len(n.Content) > 0 {
		for _, child := range n.Content {
			rendered, err := r.RenderNode(child)
			if err != nil {
				return "", err
			}
			html = append(html, rendered)
		}
	} else {
		html = append(html, template.HTMLEscapeString(n.Text))
	}

	// render closing tag(s) of node
	for i := len(tags) - 1; i >= 0; i-- {
		t := tags[i]
		closeTag, err := t.RenderClosing(n.Attrs)
		if err != nil {
			return "", fmt.Errorf("prosemirror2html: %w", err)
		}
		html = append(html, closeTag)
	}

	// render closing tags of surrounding marks
	for i := len(n.Marks) - 1; i >= 0; i-- {
		m := n.Marks[i]
		tags := r.marks[m.Type]
		if tags == nil {
			return "", fmt.Errorf("unknown mark '%s'", m.Type)
		}
		for _, t := range tags {
			closeTag, err := t.RenderClosing(m.Attrs)
			if err != nil {
				return "", fmt.Errorf("prosemirror2html: %w", err)
			}
			html = append(html, closeTag)
		}
	}

	return strings.Join(html, ""), nil
}
