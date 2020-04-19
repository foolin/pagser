package markdown

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/foolin/pagser"
	"github.com/mattn/godown"
	"regexp"
	"strings"
)

//var regMutilSpaceLine = regexp.MustCompile("(\\r?\\n\\s*){2,}")
var regMutilSpaceLine = regexp.MustCompile("(\\r?\\n+\\s*){2,}")

func Markdown(node *goquery.Selection, args ...string) (interface{}, error) {
	var buf bytes.Buffer
	html, err := node.Html()
	if err != nil {
		return "", err
	}
	err = godown.Convert(&buf, strings.NewReader(html), &godown.Option{
		Style:  true,
		Script: false,
	})
	md := buf.String()
	if err != nil {
		return md, err
	}
	md = regMutilSpaceLine.ReplaceAllString(md, "\n\n")
	return md, err
}

func Register(p *pagser.Pagser) {
	p.RegisterFunc("Markdown", Markdown)
}
