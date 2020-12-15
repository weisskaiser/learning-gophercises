package link

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// Link is used to store the information parsed from a html anchor tag (<a></a>)
// It contains the following attributes:
// - Href, used to store the href attribute if present on the tag
// - Text, used to store all the text inside the tag without any markup
type Link struct {
	Href string
	Text string
}

func isLink(node *html.Node) bool {
	return node.Type == html.ElementNode && node.Data == "a"
}

func href(node *html.Node) string {
	return attr(node, "href")
}

func attr(node *html.Node, attribute string) string {
	for _, v := range node.Attr {
		if v.Key == attribute {
			return v.Val
		}
	}
	return ""
}

func content(node *html.Node) string {

	var navigate func(node *html.Node) string

	navigate = func(node *html.Node) string {

		if node.Type == html.TextNode {
			return strings.TrimSpace(node.Data)
		}

		var content string

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			content += navigate(c) + " "
		}

		return content
	}

	content := navigate(node)

	return content[:len(content)-1]
}

func navigate(node *html.Node) ([]Link, error) {

	var links []Link

	if isLink(node) {
		l, err := newLink(node)

		if err != nil {
			return nil, err
		}

		links = append(links, *l)
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		newLinks, err := navigate(c)
		if err != nil {
			return nil, err
		}

		links = append(links, newLinks...)
	}

	return links, nil
}

func newLink(node *html.Node) (*Link, error) {

	if !isLink(node) {
		return nil, fmt.Errorf("%s Is not a link node", node.Data)
	}

	var link Link

	link.Href = href(node)
	link.Text = content(node)

	return &link, nil
}

// ParseContent is used to parse a html document to retrieve its Links
func ParseContent(htmlContent []byte) ([]Link, error) {
	node, err := html.Parse(bytes.NewReader(htmlContent))

	if err != nil {
		return nil, err
	}

	return navigate(node)
}
