package pagser

import (
	"fmt"
	"github.com/spf13/cast"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
// # Lookup function priority order
//
//	struct method -> parent method -> ... -> global
//
// # Implicit convert type
//
// 	Automatic type conversion, Output result string convert to int, int64, float64...
//
// CallFunc is a define function interface
type CallFunc func(node *goquery.Selection, args ...string) (out interface{}, err error)

// Builtin functions are registered with a lowercase initial, eg: Text -> text()
type BuiltinFunctions struct {
}

var builtinFuncObj BuiltinFunctions
var builtinFuncMap = map[string]CallFunc{
	"text":          builtinFuncObj.Text,
	"textEmpty":     builtinFuncObj.TextEmpty,
	"eachText":      builtinFuncObj.EachText,
	"eachTextEmpty": builtinFuncObj.EachTextEmpty,
	"html":          builtinFuncObj.Html,
	"eachHtml":      builtinFuncObj.EachHtml,
	"outerHtml":     builtinFuncObj.OutHtml,
	"eachOutHtml":   builtinFuncObj.EachOutHtml,
	"attr":          builtinFuncObj.Attr,
	"attrEmpty":     builtinFuncObj.AttrEmpty,
	"eachAttr":      builtinFuncObj.EachAttr,
	"eachAttrEmpty": builtinFuncObj.EachAttrEmpty,
	"attrSplit":     builtinFuncObj.AttrSplit,
	"value":         builtinFuncObj.Value,
	"split":         builtinFuncObj.Split,
	"eachJoin":      builtinFuncObj.EachJoin,
	"eq":            builtinFuncObj.Eq,
	"eqAndAttr":     builtinFuncObj.EqAndAttr,
	"eqAndHtml":     builtinFuncObj.EqAndHtml,
	"eqAndOutHtml":  builtinFuncObj.EqAndOutHtml,
	"concat":        builtinFuncObj.Concat,
	"concatAttr":    builtinFuncObj.ConcatAttr,
}

//text() get element  text, return string, this is default function, if not define function in struct tag.
//	struct {
//		Example string `pagser:".selector->text()"`
//	}
func (builtin BuiltinFunctions) Text(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return strings.TrimSpace(node.Text()), nil
}

// textEmpty(defaultValue) get element text, if empty will return defaultValue, return string.
//	struct {
//		Example string `pagser:".selector->TextEmpty('')"`
//	}
func (builtin BuiltinFunctions) TextEmpty(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("textEmpty(defaultValue) must has defaultValue")
	}
	defaultValue := args[0]
	value := strings.TrimSpace(node.Text())
	if value == "" {
		value = defaultValue
	}
	return value, nil
}

// eachText() get each element text, return []string.
//	struct {
//		Examples []string `pagser:".selector->eachText('')"`
//	}
func (builtin BuiltinFunctions) EachText(node *goquery.Selection, args ...string) (out interface{}, err error) {
	list := make([]string, 0)
	node.Each(func(i int, selection *goquery.Selection) {
		list = append(list, strings.TrimSpace(selection.Text()))
	})
	return list, nil
}

// eachTextEmpty(defaultValue) get each element text, return []string.
//	struct {
//		Examples []string `pagser:".selector->eachTextEmpty('')"`
//	}
func (builtin BuiltinFunctions) EachTextEmpty(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("eachTextEmpty(defaultValue) must has defaultValue")
	}
	defaultValue := args[0]
	list := make([]string, 0)
	node.Each(func(i int, selection *goquery.Selection) {
		value := strings.TrimSpace(selection.Text())
		if value == "" {
			value = defaultValue
		}
		list = append(list, value)
	})
	return list, nil
}

// html() get element inner html, return string.
//	struct {
//		Example string `pagser:".selector->html()"`
//	}
func (builtin BuiltinFunctions) Html(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return node.Html()
}

// eachHtml() get each element inner html, return []string.
// eachTextEmpty(defaultValue) get each element text, return []string.
//	struct {
//		Examples []string `pagser:".selector->eachHtml()"`
//	}
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
//	struct {
//		Example string `pagser:".selector->outerHtml()"`
//	}
func (builtin BuiltinFunctions) OutHtml(node *goquery.Selection, args ...string) (out interface{}, err error) {
	html, err := goquery.OuterHtml(node)
	if err != nil {
		return "", err
	}
	return html, nil
}

// eachOutHtml() get each element outer html, return []string.
//	struct {
//		Examples []string `pagser:".selector->eachOutHtml()"`
//	}
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

// attr(name, defaultValue='') get element attribute value, return string.
// outerHtml() get element  outer html, return string.
//	//<a href="https://github.com/foolin/pagser">Pagser</a>
//	struct {
//		Example string `pagser:".selector->attr(href)"`
//	}
func (builtin BuiltinFunctions) Attr(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("attr(name) must has name")
	}
	name := args[0]
	defaultValue := ""
	if len(args) > 1 {
		defaultValue = args[1]
	}
	val := node.AttrOr(name, defaultValue)
	return val, nil
}

// attrEmpty(name, defaultValue) get element attribute value, return string.
//	//<a href="https://github.com/foolin/pagser">Pagser</a>
//	struct {
//		Example string `pagser:".selector->AttrEmpty(href, '#')"`
//	}
func (builtin BuiltinFunctions) AttrEmpty(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 2 {
		return "", fmt.Errorf("attr(name, defaultValue) must has name and default value")
	}
	name := args[0]
	defaultValue := args[1]
	value := strings.TrimSpace(node.AttrOr(name, defaultValue))
	if value == "" {
		value = defaultValue
	}
	return value, nil
}

// eachAttr() get each element attribute value, return []string.
//	//<a href="https://github.com/foolin/pagser">Pagser</a>
//	struct {
//		Examples []string `pagser:".selector->eachAttr(href)"`
//	}
func (builtin BuiltinFunctions) EachAttr(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("attr(name) must has name")
	}
	name := args[0]
	list := make([]string, 0)
	node.Each(func(i int, selection *goquery.Selection) {
		list = append(list, strings.TrimSpace(selection.AttrOr(name, "")))
	})
	return list, nil
}

// eachAttrEmpty(defaultValue) get each element attribute value, return []string.
//	//<a href="https://github.com/foolin/pagser">Pagser</a>
//	struct {
//		Examples []string `pagser:".selector->eachAttrEmpty(href, '#')"`
//	}
func (builtin BuiltinFunctions) EachAttrEmpty(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 2 {
		return "", fmt.Errorf("attr(name) must has name")
	}
	name := args[0]
	defaultValue := args[1]
	list := make([]string, 0)
	node.Each(func(i int, selection *goquery.Selection) {
		value := strings.TrimSpace(selection.AttrOr(name, ""))
		if value == "" {
			value = defaultValue
		}
		list = append(list, value)
	})
	return list, nil
}

// attrSplit(name, sep=',', trim='true')  get attribute value and split by separator to array string, return []string.
//	struct {
//		Examples []string `pagser:".selector->attrSplit('keywords', ',')"`
//	}
func (builtin BuiltinFunctions) AttrSplit(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) <= 0 {
		return "", fmt.Errorf("attr(xxx) must has name")
	}
	name := args[0]
	sep := ","
	trim := true
	if len(args) > 1 {
		sep = args[1]
	}
	if len(args) > 2 {
		var err error
		trim, err = cast.ToBoolE(args[2])
		if err != nil {
			return nil, fmt.Errorf("`trim` must bool type value: true/false")
		}
	}

	list := strings.Split(node.AttrOr(name, ""), sep)
	if trim {
		for i, v := range list {
			list[i] = strings.TrimSpace(v)
		}
	}
	return list, nil
}

// value() get element attribute value by name is `value`, return string
//	//<input name="pagser" value="xxx" />
//	struct {
//		Example string `pagser:".selector->Value()"`
//	}
// Output: xxx
func (builtin BuiltinFunctions) Value(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return node.AttrOr("value", ""), nil
}

// split(sep=',', trim='true') get element text and split by separator to array string, return []string.
//	struct {
//		Examples []string `pagser:".selector->split('|')"`
//	}
func (builtin BuiltinFunctions) Split(node *goquery.Selection, args ...string) (out interface{}, err error) {
	sep := ","
	trim := true
	if len(args) > 0 {
		sep = args[0]
	}
	if len(args) > 1 {
		var err error
		trim, err = cast.ToBoolE(args[1])
		if err != nil {
			return nil, fmt.Errorf("`trim` must bool type value: true/false")
		}
	}
	list := strings.Split(node.Text(), sep)
	if trim {
		for i, v := range list {
			list[i] = strings.TrimSpace(v)
		}
	}
	return list, nil
}

// eachJoin(sep) get each element text and join to string, return string.
//	struct {
//		Example string `pagser:".selector->eachJoin(',')"`
//	}
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
//	struct {
//		Example string `pagser:".selector->eq(0)"`
//	}
func (builtin BuiltinFunctions) Eq(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) <= 0 {
		return "", fmt.Errorf("eq(index) must has index")
	}
	indexValue := strings.TrimSpace(args[0])
	idx, err := strconv.Atoi(indexValue)
	if err != nil {
		return "", fmt.Errorf("index=`" + indexValue + "` is not number: " + err.Error())
	}
	return strings.TrimSpace(node.Eq(idx).Text()), nil
}

// eqAndAttr(index, name) reduces the set of matched elements to the one at the specified index, and attr() return string.
//	struct {
//		Example string `pagser:".selector->eqAndAttr(0, href)"`
//	}
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
//	struct {
//		Example string `pagser:".selector->eqAndHtml(0)"`
//	}
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
//	struct {
//		Example string `pagser:".selector->eqAndOutHtml(0)"`
//	}
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

// concat(text1, $value, [ text2, ... text_n ])
// The `text1, text2, ... text_n` strings that you wish to join together,
// `$value` is placeholder for get element  text, return string.
//	struct {
//		Example string `pagser:".selector->concat('Result:', '<', $value, '>')"`
//	}
func (builtin BuiltinFunctions) Concat(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 2 {
		return "", fmt.Errorf("concat(text1, $value, [ text2, ... text_n ]) must be more than two arguments")
	}
	value := strings.TrimSpace(node.Text())
	builder := strings.Builder{}
	for _, v := range args {
		if v == "$value" {
			builder.WriteString(value)
		} else {
			builder.WriteString(v)
		}
	}
	return builder.String(), nil
}

// concatAttr(name, text1, $value, [ text2, ... text_n ])
// `name` get element attribute value by name,
// `text1, text2, ... text_n` The strings that you wish to join together,
// `$value` is placeholder for get element  text
// return string.
//	struct {
//		Example string `pagser:".selector->concatAttr('Result:', '<', $value, '>')"`
//	}
func (builtin BuiltinFunctions) ConcatAttr(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 3 {
		return "", fmt.Errorf("concatAttr(name, text1, $value, [ text2, ... text_n ]) must be more than two arguments")
	}
	name := args[0]
	value := strings.TrimSpace(node.AttrOr(name, ""))
	builder := strings.Builder{}
	for i, v := range args {
		if i == 0 {
			continue
		}
		if v == "$value" {
			builder.WriteString(value)
		} else {
			builder.WriteString(v)
		}
	}
	return builder.String(), nil
}

// RegisterFunc register function for parse result
//	pagser.RegisterFunc("MyFunc", func(node *goquery.Selection, args ...string) (out interface{}, err error) {
//		//Todo
//		return "Hello", nil
//	})
func (p *Pagser) RegisterFunc(name string, fn CallFunc) {
	p.ctxFuncs[name] = fn
}
