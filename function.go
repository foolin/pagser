package pagser

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mattn/godown"
	"github.com/microcosm-cc/bluemonday"
	"regexp"
	"strings"
)

//var regMutilSpaceLine = regexp.MustCompile("(\\r?\\n\\s*){2,}")
var regMutilSpaceLine = regexp.MustCompile("(\\r?\\n+\\s*){2,}")

type CallFunc func(node *goquery.Selection, args ...string) (out interface{}, err error)

var sysFuncs = map[string]CallFunc{
	"html":      html,
	"text":      text,
	"attr":      attr,
	"markdown":  markdown,
	"ugcHtml":   ugcHtml,
	"outerHtml": outHtml,
	"value":     value,
}

func markdown(node *goquery.Selection, args ...string) (interface{}, error) {
	var buf bytes.Buffer
	html, err := ugcHtml(node)
	if err != nil {
		return "", err
	}
	err = godown.Convert(&buf, strings.NewReader(html.(string)), &godown.Option{
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

func ugcHtml(node *goquery.Selection, args ...string) (interface{}, error) {
	html, err := goquery.OuterHtml(node)
	if err != nil {
		return html, err
	}
	p := bluemonday.UGCPolicy()
	// The policy can then be used to sanitize lots of input and it is safe to use the policy in multiple goroutines
	return p.Sanitize(html), nil
}

func html(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return node.Html()
}

func outHtml(node *goquery.Selection, args ...string) (out interface{}, err error) {
	html, err := goquery.OuterHtml(node)
	if err != nil {
		return "", err
	}
	return html, nil
}

func value(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return node.AttrOr("value", ""), nil
}

func text(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return strings.TrimSpace(node.Text()), nil
}

func attr(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) <= 0 {
		return "", fmt.Errorf("attr(xxx) must has name")
	}
	val, _ := node.Attr(args[0])
	return val, nil
}

func (p *Pagser) RegisterFunc(name string, fn CallFunc) error {
	p.funcs[name] = fn
	return nil
}
