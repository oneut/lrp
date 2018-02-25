package proxy

import (
	"bytes"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
)

type ProxyDocument struct {
	Body io.ReadCloser
}

func (pd *ProxyDocument) CreateBytesBufferWithLiveReloadScriptPath(scriptPath string) bytes.Buffer {
	defer pd.Body.Close()
	document, err := html.Parse(pd.Body)
	if err != nil {
		panic(err)
	}
	body := pd.FindFirstChild(document, atom.Body)
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
	return pd.ConvertBytesBufferFromHTMLNode(document)
}

func (pd *ProxyDocument) FindFirstChild(node *html.Node, a atom.Atom) *html.Node {
	for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		if childNode.Type != html.ElementNode {
			continue
		}
		if childNode.DataAtom == a {
			return node
		}
		if n := pd.FindFirstChild(childNode, a); n != nil {
			return n
		}

	}

	return nil
}

func (pd *ProxyDocument) ConvertBytesBufferFromHTMLNode(node *html.Node) bytes.Buffer {
	var buf bytes.Buffer
	for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		html.Render(&buf, childNode)
	}

	return buf
}
