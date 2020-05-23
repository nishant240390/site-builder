package main

import(
	"io"
	"strings"
	"golang.org/x/net/html"
)


type Link struct {
	Href string
	Text string
}
func Parse(r io.Reader)([]Link,error) {
	root ,err := html.Parse(r)

	if err != nil {
		return nil, err
	}
	var links []Link

	var rec func(*html.Node)

	rec = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					var text string
					if n.FirstChild != nil {
						text = grabText(n.FirstChild)
					}
					links = append(links,Link{attr.Val,text})
				}
			}
		}
		if n.FirstChild != nil {
			rec(n.FirstChild)
		}
		if n.NextSibling != nil {
			rec (n.NextSibling)
		}
	}
	rec (root)

	return links, nil
}

func grabText (n *html.Node) string {
	var sb strings.Builder
	var rec func(*html.Node)
	rec = func (n *html.Node) {
		if n.Type == html.TextNode {
			s := n.Data
			sb.WriteString(s)
		}
		if n.FirstChild != nil {
			rec(n.FirstChild)
		}
		if n.NextSibling != nil {
			rec(n.NextSibling)
		}
	}
	rec(n)
	return strings.Join(strings.Fields(sb.String()), " ")
}
