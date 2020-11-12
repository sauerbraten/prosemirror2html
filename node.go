package prosemirror2html

// Node is a Go representation of http://prosemirror.net/docs/ref/#model.Node.
type Node struct {
	Type    string                 `json:"type"`
	Attrs   map[string]interface{} `json:"attrs,omitempty"`
	Content []*Node                `json:"content,omitempty"`
	Marks   []mark                 `json:"marks,omitempty"`
	Text    string                 `json:"text,omitempty"`
}

// http://prosemirror.net/docs/ref/#model.Mark
type mark struct {
	Type  string                 `json:"type"`
	Attrs map[string]interface{} `json:"attrs,omitempty"`
}
