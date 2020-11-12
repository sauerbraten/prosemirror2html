package prosemirror2html

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestRenderer(t *testing.T) {
	testcases := []struct {
		input    string
		expected string
	}{
		{
			input:    `{"type":"doc"}`,
			expected: ``,
		},
		{
			input:    `{"type":"doc","content":[{"type":"text"}]}`,
			expected: ``,
		},
		{
			input:    `{"type":"doc","content":[{"type":"text","text":""}]}`,
			expected: ``,
		},
		{
			input:    `{"type":"doc","content":[{"type":"text","text":null}]}`,
			expected: ``,
		},
		{
			input:    `{"type":"doc","content":[{"type":"text","text":"foo bar"}]}`,
			expected: `foo bar`,
		},
		{
			input:    `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"foo bar "},{"type":"text","text":"asd","marks":[{"type":"link","attrs":{"href":"https://www.spiegel.de/","target":"_blank","rel":"noopener noreferrer"}}]}]}]}`,
			expected: `<p>foo bar <a href="https://www.spiegel.de/" target="_blank" rel="noopener noreferrer">asd</a></p>`,
		},
		{
			input:    `{"type":"doc","content":[{"type":"heading","attrs":{"level":2},"content":[{"type":"text","text":"Example Heading"}]},{"type":"paragraph","content":[{"type":"text","text":"You can "},{"type":"text","marks":[{"type":"code"}],"text":"write code"},{"type":"text","text":"."}]},{"type":"paragraph","content":[{"type":"text","text":"There are lots of formatting options, like "},{"type":"text","marks":[{"type":"bold"}],"text":"bold"},{"type":"text","text":", "},{"type":"text","marks":[{"type":"italic"}],"text":"italics"},{"type":"text","text":", "},{"type":"text","marks":[{"type":"underline"}],"text":"underline"},{"type":"text","text":", and "},{"type":"text","marks":[{"type":"strike"}],"text":"strikethrough"},{"type":"text","text":"."}]},{"type":"bullet_list","content":[{"type":"list_item","content":[{"type":"paragraph","content":[{"type":"text","text":"there are"}]}]},{"type":"list_item","content":[{"type":"paragraph","content":[{"type":"text","text":"bullet lists"}]}]}]},{"type":"paragraph","content":[{"type":"text","text":"as well as"}]},{"type":"ordered_list","attrs":{"order":1},"content":[{"type":"list_item","content":[{"type":"paragraph","content":[{"type":"text","text":"ordered"}]}]},{"type":"list_item","content":[{"type":"paragraph","content":[{"type":"text","text":"lists"}]}]}]},{"type":"blockquote","content":[{"type":"paragraph","content":[{"type":"text","text":"You can also make blockquotes"}]}]}]  }`,
			expected: `<h2>Example Heading</h2><p>You can <code>write code</code>.</p><p>There are lots of formatting options, like <strong>bold</strong>, <em>italics</em>, <u>underline</u>, and <s>strikethrough</s>.</p><ul><li><p>there are</p></li><li><p>bullet lists</p></li></ul><p>as well as</p><ol order=1><li><p>ordered</p></li><li><p>lists</p></li></ol><blockquote><p>You can also make blockquotes</p></blockquote>`,
		},
	}

	r := NewRenderer()

	for _, c := range testcases {
		expectedNodes, err := html.ParseFragment(strings.NewReader(c.expected), &html.Node{
			Type:     html.ElementNode,
			Data:     "body",
			DataAtom: atom.Body,
		})
		if err != nil {
			t.Fatal("could not parse expected HTML:", c.expected)
		}

		output, err := r.Render([]byte(c.input))
		if err != nil {
			t.Fatal(err)
		}

		outputNodes, err := html.ParseFragment(strings.NewReader(output), &html.Node{
			Type:     html.ElementNode,
			Data:     "body",
			DataAtom: atom.Body,
		})
		if err != nil {
			t.Fatal("could not parse output HTML:", output)
		}

		if len(outputNodes) != len(expectedNodes) {
			t.Fatal("\ngot:\n\t", output, "\nbut expected:\n\t", c.expected)
		}
		for i := range outputNodes {
			if !equal(outputNodes[i], expectedNodes[i]) {
				t.Fatal("\ngot:\n\t", output, "\nbut expected:\n\t", c.expected)
			}
		}
	}
}

func equal(one, two *html.Node) bool {
	if one == nil || two == nil {
		return one == two
	}

	if one.Type != two.Type || one.Data != two.Data || len(one.Attr) != len(two.Attr) {
		return false
	}

outer:
	for _, attrInOne := range one.Attr {
		for _, attrInTwo := range two.Attr {
			if attrInOne.Key == attrInTwo.Key && attrInOne.Val == attrInTwo.Val {
				continue outer
			}
		}
		return false
	}

	return equal(one.NextSibling, two.NextSibling) &&
		equal(one.FirstChild, two.FirstChild)
}
