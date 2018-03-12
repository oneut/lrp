package livereloadproxy

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type ProxyBody struct {
	Body io.ReadCloser
}

func (pb *ProxyBody) getBytesBufferWithLiveReloadScriptPath(scriptPath string) bytes.Buffer {
	defer pb.Body.Close()
	document, err := html.Parse(pb.Body)
	if err != nil {
		panic(err)
	}

	body := pb.FindFirstChild(document, atom.Body)
	body.AppendChild(&html.Node{
		Type:     html.ElementNode,
		Data:     "script",
		DataAtom: atom.Script,
		Attr: []html.Attribute{
			html.Attribute{
				Key: "src",
				Val: scriptPath,
			},
		},
	})

	var buf bytes.Buffer
	html.Render(&buf, document)
	return buf
}

func (pb *ProxyBody) FindFirstChild(node *html.Node, a atom.Atom) *html.Node {
	for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		if childNode.Type != html.ElementNode {
			continue
		}

		if childNode.DataAtom == a {
			return childNode
		}
		if n := pb.FindFirstChild(childNode, a); n != nil {
			return n
		}

	}

	return nil
}
