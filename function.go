package pagser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

// CallFunc write function interface
//
// Builtin Functions
//
// - text() get element  text, return string, this is default function, if not define function in struct tag.
//
// - eachText() get each element text, return []string.
//
// - html() get element inner html, return string.
//
// - eachHtml() get each element inner html, return []string.
//
// - outerHtml() get element  outer html, return string.
//
// - eachOutHtml() get each element outer html, return []string.
//
// - attr(name) get element attribute value, return string.
//
// - eachAttr() get each element attribute value, return []string.
//
// - attrInt(name, defaultValue) get element attribute value and to int, return int.
//
// - attrSplit(name, sep)  get attribute value and split by separator to array string.
//
// - value() get element attribute value by name is `value`, return string, eg: <input value='xxxx' /> will return "xxx".
//
// - split(sep) get element text and split by separator to array string, return []string.
//
// - eachJoin(sep) get each element text and join to string, return string.
//
//
// # Define Global Function
//
//	func MyFunc(node *goquery.Selection, args ...string) (out interface{}, err error) {
//		//Todo
//		return "Hello", nil
//	}
//
//	//Register function
//	pagser.RegisterFunc("MyFunc", MyFunc)
//
//	//Use function
//	type PageData struct{
//	     Text string `pagser:"h1->MyFunc()"`
//	}
//
//
// # Define Struct Function
//	//Use function
//	type PageData struct{
//	     Text string `pagser:"h1->MyFunc()"`
//	}
//
//	func (my PageData) MyFunc(node *goquery.Selection, args ...string) (out interface{}, err error) {
//		//Todo
//		return "Hello", nil
//	}
//
//
//
// Define your own function interface
type CallFunc func(node *goquery.Selection, args ...string) (out interface{}, err error)

var sysFuncs = map[string]CallFunc{
	"text":        text,
	"eachText":    eachText,
	"html":        html,
	"eachHtml":    eachHtml,
	"outerHtml":   outHtml,
	"eachOutHtml": eachOutHtml, //
	"attr":        attr,        //
	"eachAttr":    eachAttr,
	"attrInt":     attrInt,
	"attrSplit":   attrSplit,
	"value":       value,
	"split":       split,
	"eachJoin":    eachJoin,
}

// text() string
func text(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return strings.TrimSpace(node.Text()), nil
}

// eachText() []string
func eachText(node *goquery.Selection, args ...string) (out interface{}, err error) {
	list := make([]string, 0)
	node.Each(func(i int, selection *goquery.Selection) {
		list = append(list, strings.TrimSpace(selection.Text()))
	})
	return list, nil
}

// html() string
func html(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return node.Html()
}

// eachHtml() []string
func eachHtml(node *goquery.Selection, args ...string) (out interface{}, err error) {
	list := make([]string, 0)
	node.EachWithBreak(func(i int, selection *goquery.Selection) bool {
		var html string
		html, err = node.Html()
		if err != nil {
			return false
		}
		list = append(list, html)
		return true
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// outHtml() string
func outHtml(node *goquery.Selection, args ...string) (out interface{}, err error) {
	html, err := goquery.OuterHtml(node)
	if err != nil {
		return "", err
	}
	return html, nil
}

// eachOutHtml() []string
func eachOutHtml(node *goquery.Selection, args ...string) (out interface{}, err error) {
	list := make([]string, 0)
	node.EachWithBreak(func(i int, selection *goquery.Selection) bool {
		var html string
		html, err = goquery.OuterHtml(node)
		if err != nil {
			return false
		}
		list = append(list, html)
		return true
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// attr(name) string
func attr(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) <= 0 {
		return "", fmt.Errorf("attr(xxx) must has name")
	}
	name := args[0]
	val, _ := node.Attr(name)
	return val, nil
}

// eachAttr(name) []string
func eachAttr(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) <= 0 {
		return "", fmt.Errorf("attr(xxx) must has name")
	}
	name := args[0]
	list := make([]string, 0)
	node.Each(func(i int, selection *goquery.Selection) {
		list = append(list, strings.TrimSpace(selection.AttrOr(name, "")))
	})
	return list, nil
}

// attrInt(name) int
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

// attrSplit(name, sep) []string
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

// value() string
func value(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return node.AttrOr("value", ""), nil
}

// split(sep) []string
func split(node *goquery.Selection, args ...string) (out interface{}, err error) {
	sep := ","
	if len(args) > 0 {
		sep = args[0]
	}
	return strings.Split(node.Text(), sep), nil
}

// eachJoin(sep) string
func eachJoin(node *goquery.Selection, args ...string) (out interface{}, err error) {
	sep := ","
	if len(args) > 0 {
		sep = args[0]
	}
	list := make([]string, 0)
	node.Each(func(i int, selection *goquery.Selection) {
		list = append(list, strings.TrimSpace(selection.Text()))
	})
	return strings.Join(list, sep), nil
}

// RegisterFunc register function for parse
func (p *Pagser) RegisterFunc(name string, fn CallFunc) error {
	p.funcs[name] = fn
	return nil
}
