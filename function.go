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
	"attr":          builtinFuncObj.Attr,
	"attrConcat":    builtinFuncObj.AttrConcat,
	"attrEmpty":     builtinFuncObj.AttrEmpty,
	"attrSplit":     builtinFuncObj.AttrSplit,
	"eachAttr":      builtinFuncObj.EachAttr,
	"eachAttrEmpty": builtinFuncObj.EachAttrEmpty,
	"eachHtml":      builtinFuncObj.EachHtml,
	"eachOutHtml":   builtinFuncObj.EachOutHtml,
	"eachText":      builtinFuncObj.EachText,
	"eachTextEmpty": builtinFuncObj.EachTextEmpty,
	"eachTextJoin":  builtinFuncObj.EachTextJoin,
	"eqAndAttr":     builtinFuncObj.EqAndAttr,
	"eqAndHtml":     builtinFuncObj.EqAndHtml,
	"eqAndOutHtml":  builtinFuncObj.EqAndOutHtml,
	"eqAndText":     builtinFuncObj.EqAndText,
	"html":          builtinFuncObj.Html,
	"nodeChild":     builtinFuncObj.NodeChild,
	"nodeEq":        builtinFuncObj.NodeEq,
	"nodeNext":      builtinFuncObj.NodeNext,
	"nodeParent":    builtinFuncObj.NodeParent,
	"nodePrev":      builtinFuncObj.NodePrev,
	"nodeSiblings":  builtinFuncObj.NodeSiblings,
	"outerHtml":     builtinFuncObj.OutHtml,
	"text":          builtinFuncObj.Text,
	"textConcat":    builtinFuncObj.TextConcat,
	"textEmpty":     builtinFuncObj.TextEmpty,
	"textSplit":     builtinFuncObj.TextSplit,
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

// attrConcat(name, text1, $value, [ text2, ... text_n ])
// `name` get element attribute value by name,
// `text1, text2, ... text_n` The strings that you wish to join together,
// `$value` is placeholder for get element  text
// return string.
//	struct {
//		Example string `pagser:".selector->attrConcat('Result:', '<', $value, '>')"`
//	}
func (builtin BuiltinFunctions) AttrConcat(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 3 {
		return "", fmt.Errorf("attrConcat(name, text1, $value, [ text2, ... text_n ]) must be more than two arguments")
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

// attrSplit(name, sep=',', trim='true')  get attribute value and split by separator to array string, return []string.
//	struct {
//		Examples []string `pagser:".selector->attrSplit('keywords', ',')"`
//	}
func (builtin BuiltinFunctions) AttrSplit(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("attr(name) must has `name`")
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

// eachTextJoin(sep) get each element text and join to string, return string.
//	struct {
//		Example string `pagser:".selector->eachTextJoin(',')"`
//	}
func (builtin BuiltinFunctions) EachTextJoin(node *goquery.Selection, args ...string) (out interface{}, err error) {
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

// eqAndAttr(index, name) reduces the set of matched elements to the one at the specified index, and attr() return string.
//	struct {
//		Example string `pagser:".selector->eqAndAttr(0, href)"`
//	}
func (builtin BuiltinFunctions) EqAndAttr(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 2 {
		return "", fmt.Errorf("eqAndAttr(index) must has index and attr name")
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
	if len(args) < 1 {
		return "", fmt.Errorf("eqAndHtml(index) must has index")
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
	if len(args) < 1 {
		return "", fmt.Errorf("eqAndOutHtml(index) must has index")
	}
	indexValue := strings.TrimSpace(args[0])
	idx, err := strconv.Atoi(indexValue)
	if err != nil {
		return "", fmt.Errorf("index=`" + indexValue + "` is not number: " + err.Error())
	}
	return goquery.OuterHtml(node.Eq(idx))
}

// eqAndText(index) reduces the set of matched elements to the one at the specified index, return string.
//	struct {
//		Example string `pagser:".selector->eqAndText(0)"`
//	}
func (builtin BuiltinFunctions) EqAndText(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("eqAndText(index) must has index")
	}
	indexValue := strings.TrimSpace(args[0])
	idx, err := strconv.Atoi(indexValue)
	if err != nil {
		return "", fmt.Errorf("index=`" + indexValue + "` is not number: " + err.Error())
	}
	return strings.TrimSpace(node.Eq(idx).Text()), nil
}

// html() get element inner html, return string.
//	struct {
//		Example string `pagser:".selector->html()"`
//	}
func (builtin BuiltinFunctions) Html(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return node.Html()
}

// nodeChild(selector = '') gets the child elements of each element in the Selection,
// Filtered by the specified selector if selector not empty,
// It returns Selection object containing these elements.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->nodeChild()"`
//	}
func (builtin BuiltinFunctions) NodeChild(node *goquery.Selection, args ...string) (out interface{}, err error) {
	selector := ""
	if len(args) > 0 {
		selector = strings.TrimSpace(args[0])
	}
	if selector != "" {
		return node.ChildrenFiltered(selector), nil
	}
	return node.Children(), nil
}

// nodeEq(index) reduces the set of matched elements to the one at the specified index.
// If a negative index is given, it counts backwards starting at the end of the
// set. It returns a Selection object, and an empty Selection object if the
// index is invalid.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->nodeEq(0)"`
//	}
func (builtin BuiltinFunctions) NodeEq(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("nodeEq(index) must has `index` value")
	}
	indexValue := strings.TrimSpace(args[0])
	idx, err := strconv.Atoi(indexValue)
	if err != nil {
		return "", fmt.Errorf("index=`" + indexValue + "` is not number: " + err.Error())
	}
	return node.Eq(idx), nil
}

// nodeNext() gets the immediately following sibling of each element in the Selection.
// Filtered by the specified selector if selector not empty,
// It returns Selection object containing these elements.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->nodeNext()"`
//	}
func (builtin BuiltinFunctions) NodeNext(node *goquery.Selection, args ...string) (out interface{}, err error) {
	selector := ""
	if len(args) > 0 {
		selector = strings.TrimSpace(args[0])
	}
	if selector != "" {
		return node.NextFiltered(selector), nil
	}
	return node.Next(), nil
}

// nodeParent() gets the parent elements of each element in the Selection.
// Filtered by the specified selector if selector not empty,
// It returns Selection object containing these elements.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->nodeParent()"`
//	}
func (builtin BuiltinFunctions) NodeParent(node *goquery.Selection, args ...string) (out interface{}, err error) {
	selector := ""
	if len(args) > 0 {
		selector = strings.TrimSpace(args[0])
	}
	if selector != "" {
		return node.ParentFiltered(selector), nil
	}
	return node.Parent(), nil
}

// nodePrev() gets the immediately preceding sibling of each element in the Selection.
// Filtered by the specified selector if selector not empty,
// It returns Selection object containing these elements.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->nodePrev()"`
//	}
func (builtin BuiltinFunctions) NodePrev(node *goquery.Selection, args ...string) (out interface{}, err error) {
	selector := ""
	if len(args) > 0 {
		selector = strings.TrimSpace(args[0])
	}
	if selector != "" {
		return node.PrevFiltered(selector), nil
	}
	return node.Prev(), nil
}

// nodeSiblings() gets the siblings of each element in the Selection.
// Filtered by the specified selector if selector not empty,
// It returns Selection object containing these elements.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->nodeSiblings()"`
//	}
func (builtin BuiltinFunctions) NodeSiblings(node *goquery.Selection, args ...string) (out interface{}, err error) {
	selector := ""
	if len(args) > 0 {
		selector = strings.TrimSpace(args[0])
	}
	if selector != "" {
		return node.SiblingsFiltered(selector), nil
	}
	return node.Siblings(), nil
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

//text() get element  text, return string, this is default function, if not define function in struct tag.
//	struct {
//		Example string `pagser:".selector->text()"`
//	}
func (builtin BuiltinFunctions) Text(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return strings.TrimSpace(node.Text()), nil
}

// textConcat(text1, $value, [ text2, ... text_n ])
// The `text1, text2, ... text_n` strings that you wish to join together,
// `$value` is placeholder for get element  text, return string.
//	struct {
//		Example string `pagser:".selector->textConcat('Result:', '<', $value, '>')"`
//	}
func (builtin BuiltinFunctions) TextConcat(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 2 {
		return "", fmt.Errorf("textConcat(text1, $value, [ text2, ... text_n ]) must be more than two arguments")
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

// textSplit(sep=',', trim='true') get element text and split by separator to array string, return []string.
//	struct {
//		Examples []string `pagser:".selector->textSplit('|')"`
//	}
func (builtin BuiltinFunctions) TextSplit(node *goquery.Selection, args ...string) (out interface{}, err error) {
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

// RegisterFunc register function for parse result
//	pagser.RegisterFunc("MyFunc", func(node *goquery.Selection, args ...string) (out interface{}, err error) {
//		//Todo
//		return "Hello", nil
//	})
func (p *Pagser) RegisterFunc(name string, fn CallFunc) {
	p.ctxFuncs[name] = fn
}
