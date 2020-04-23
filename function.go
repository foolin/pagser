package pagser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

// CallFunc write function interface

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
//	func (pd PageData) MyFunc(node *goquery.Selection, args ...string) (out interface{}, err error) {
//		//Todo
//		return "Hello", nil
//	}
//
//
//
// Define your own function interface
type CallFunc func(node *goquery.Selection, args ...string) (out interface{}, err error)

// Builtin functions are registered with a lowercase initial, eg: Text -> text()
type BuiltinFunctions struct {
}

var builtinFuncObj BuiltinFunctions
var builtinFuncMap = map[string]CallFunc{
	"text":         builtinFuncObj.Text,
	"eachText":     builtinFuncObj.EachText,
	"html":         builtinFuncObj.Html,
	"eachHtml":     builtinFuncObj.EachHtml,
	"outerHtml":    builtinFuncObj.OutHtml,
	"eachOutHtml":  builtinFuncObj.EachOutHtml, //
	"attr":         builtinFuncObj.Attr,        //
	"eachAttr":     builtinFuncObj.EachAttr,
	"attrInt":      builtinFuncObj.AttrInt,
	"attrSplit":    builtinFuncObj.AttrSplit,
	"value":        builtinFuncObj.Value,
	"split":        builtinFuncObj.Split,
	"eachJoin":     builtinFuncObj.EachJoin,
	"eq":           builtinFuncObj.Eq,
	"eqAndAttr":    builtinFuncObj.EqAndAttr,
	"eqAndHtml":    builtinFuncObj.EqAndHtml,
	"eqAndOutHtml": builtinFuncObj.EqAndOutHtml,
}

// text() get element  text, return string, this is default function, if not define function in struct tag.
func (builtin BuiltinFunctions) Text(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return strings.TrimSpace(node.Text()), nil
}

// eachText() get each element text, return []string.
func (builtin BuiltinFunctions) EachText(node *goquery.Selection, args ...string) (out interface{}, err error) {
	list := make([]string, 0)
	node.Each(func(i int, selection *goquery.Selection) {
		list = append(list, strings.TrimSpace(selection.Text()))
	})
	return list, nil
}

// html() get element inner html, return string.
func (builtin BuiltinFunctions) Html(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return node.Html()
}

// eachHtml() get each element inner html, return []string.
func (builtin BuiltinFunctions) EachHtml(node *goquery.Selection, args ...string) (out interface{}, err error) {
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

// outerHtml() get element  outer html, return string.
func (builtin BuiltinFunctions) OutHtml(node *goquery.Selection, args ...string) (out interface{}, err error) {
	html, err := goquery.OuterHtml(node)
	if err != nil {
		return "", err
	}
	return html, nil
}

// eachOutHtml() get each element outer html, return []string.
func (builtin BuiltinFunctions) EachOutHtml(node *goquery.Selection, args ...string) (out interface{}, err error) {
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

// attr(name) get element attribute value, return string.
func (builtin BuiltinFunctions) Attr(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) <= 0 {
		return "", fmt.Errorf("attr(xxx) must has name")
	}
	name := args[0]
	val, _ := node.Attr(name)
	return val, nil
}

// eachAttr() get each element attribute value, return []string.
func (builtin BuiltinFunctions) EachAttr(node *goquery.Selection, args ...string) (out interface{}, err error) {
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

// attrInt(name, defaultValue) get element attribute value and to int, return int.
func (builtin BuiltinFunctions) AttrInt(node *goquery.Selection, args ...string) (out interface{}, err error) {
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

// attrSplit(name, sep)  get attribute value and split by separator to array string.
func (builtin BuiltinFunctions) AttrSplit(node *goquery.Selection, args ...string) (out interface{}, err error) {
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

// value() get element attribute value by name is `value`, return string
func (builtin BuiltinFunctions) Value(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return node.AttrOr("value", ""), nil
}

// split(sep) get element text and split by separator to array string, return []string.
func (builtin BuiltinFunctions) Split(node *goquery.Selection, args ...string) (out interface{}, err error) {
	sep := ","
	if len(args) > 0 {
		sep = args[0]
	}
	return strings.Split(node.Text(), sep), nil
}

// eachJoin(sep) get each element text and join to string, return string.
func (builtin BuiltinFunctions) EachJoin(node *goquery.Selection, args ...string) (out interface{}, err error) {
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

// eq(index) reduces the set of matched elements to the one at the specified index, return string.
func (builtin BuiltinFunctions) Eq(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) <= 0 {
		return "", fmt.Errorf("eq(index) must has index")
	}
	indexValue := strings.TrimSpace(args[0])
	idx, err := strconv.Atoi(indexValue)
	if err != nil {
		return "", fmt.Errorf("index=`" + indexValue + "` is not number: " + err.Error())
	}
	return node.Eq(idx).Text(), nil
}

// eqAndAttr(index, name) reduces the set of matched elements to the one at the specified index, and attr() return string.
func (builtin BuiltinFunctions) EqAndAttr(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) <= 1 {
		return "", fmt.Errorf("eq(index) must has index and attr name")
	}
	indexValue := strings.TrimSpace(args[0])
	idx, err := strconv.Atoi(indexValue)
	if err != nil {
		return "", fmt.Errorf("index=`" + indexValue + "` is not number: " + err.Error())
	}
	name := strings.TrimSpace(args[1])
	return node.Eq(idx).AttrOr(name, ""), nil
}

// eqAndHtml(index) reduces the set of matched elements to the one at the specified index, and html() return string.
func (builtin BuiltinFunctions) EqAndHtml(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) <= 0 {
		return "", fmt.Errorf("eq(index) must has index")
	}
	indexValue := strings.TrimSpace(args[0])
	idx, err := strconv.Atoi(indexValue)
	if err != nil {
		return "", fmt.Errorf("index=`" + indexValue + "` is not number: " + err.Error())
	}
	return node.Eq(idx).Html()
}

// eqAndOutHtml(index) reduces the set of matched elements to the one at the specified index, and outHtml() return string.
func (builtin BuiltinFunctions) EqAndOutHtml(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) <= 0 {
		return "", fmt.Errorf("eq(index) must has index")
	}
	indexValue := strings.TrimSpace(args[0])
	idx, err := strconv.Atoi(indexValue)
	if err != nil {
		return "", fmt.Errorf("index=`" + indexValue + "` is not number: " + err.Error())
	}
	return goquery.OuterHtml(node.Eq(idx))
}

// RegisterFunc register function for parse
func (p *Pagser) RegisterFunc(name string, fn CallFunc) error {
	p.ctxFuncs[name] = fn
	return nil
}
