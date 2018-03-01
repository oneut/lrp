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

func (pb *ProxyBody) CreateBytesBufferWithLiveReloadScriptPath(scriptPath string) bytes.Buffer {
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
	return pb.ConvertBytesBufferFromHTMLNode(document)
}

func (pb *ProxyBody) FindFirstChild(node *html.Node, a atom.Atom) *html.Node {
	for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		if childNode.Type != html.ElementNode {
			continue
		}
		if childNode.DataAtom == a {
			return node
		}
		if n := pb.FindFirstChild(childNode, a); n != nil {
			return n
		}

	}

	return nil
}

func (pb *ProxyBody) ConvertBytesBufferFromHTMLNode(node *html.Node) bytes.Buffer {
	var buf bytes.Buffer
	for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		html.Render(&buf, childNode)
	}

	return buf
}
