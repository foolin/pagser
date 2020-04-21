package pagser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

// CallFunc write function interface
type CallFunc func(node *goquery.Selection, args ...string) (out interface{}, err error)

var sysFuncs = map[string]CallFunc{
	"html":      html,
	"text":      text,
	"attr":      attr,
	"attrInt":   attrInt,
	"outerHtml": outHtml,
	"value":     value,
	"split":     split,
	"attrSplit": attrSplit,
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

func attrSplit(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) <= 0 {
		return "", fmt.Errorf("attr(xxx) must has name")
	}
	name := args[0]
	sep := ","
	if len(args) > 1 {
		sep = args[1]
	}
	return strings.Split(node.AttrOr(name, ""), sep), nil
}

func attrInt(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 2 {
		return "", fmt.Errorf("attrInt(name,defaultValue) must has name and default value, eg: attrInt(id,-1)")
	}
	name := args[0]
	defaultValue := args[1]
	attrVal := node.AttrOr(name, defaultValue)
	outVal, err := strconv.Atoi(attrVal)
	if err != nil {
		return strconv.Atoi(defaultValue)
	}
	return outVal, nil
}

// split split
func split(node *goquery.Selection, args ...string) (out interface{}, err error) {
	sep := ","
	if len(args) > 0 {
		sep = args[0]
	}
	return strings.Split(node.Text(), sep), nil
}

// RegisterFunc register function for parse
func (p *Pagser) RegisterFunc(name string, fn CallFunc) error {
	p.funcs[name] = fn
	return nil
}
