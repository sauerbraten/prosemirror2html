# prosemirror2html

This is a Go reimplementation of https://github.com/ueberdosis/prosemirror-to-html, supporting all the same nodes and marks by default.

Documentation is at https://pkg.go.dev/github.com/sauerbraten/prosemirror2html.

## Extending

To use this package with custom node or mark types, simply implement the `Tag` interface and register your type using `RegisterMark()` or `RegisterNode()`. The provided `SimpleTag` type is often enough:

```go
package main

import (
	"github.com/sauerbraten/prosemirror2html"
)

type CustomMark struct{}

func (m CustomMark) RenderOpening(attrs map[string]interface{}) (string, error) {
	return "<custom-opening-tag " + attrs["my_special_attribute"].(string) + ">", nil
}

func (m CustomMark) RenderClosing(map[string]interface{}) (string, error) {
	return "</custom-closing-tag>", nil
}

func main() {
	r := prosemirror2html.NewRenderer()
	r.RegisterNode("my_node_type", prosemirror2html.SimpleTag{Name: "my-node-type"})
    r.RegisterMark("my_mark_type", CustomMark{})

    // ...
}
```