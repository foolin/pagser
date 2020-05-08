package pagser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cast"
	"net/url"
	"strconv"
	"strings"
)

// BuiltinFunctions builtin functions are registered with a lowercase initial, eg: Text -> text()
type BuiltinFunctions struct {
}

// AbsHref absHref(baseUrl) get element attribute name `href`, and convert to absolute url, return *URL.
// `baseUrl` is the base url like `https://example.com/`.
//	//<a href="/foolin/pagser">Pagser</a>
//	struct {
//		Example string `pagser:".selector->absHref('https://github.com/')"`
//	}
func (builtin BuiltinFunctions) AbsHref(selection *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("args must has baseUrl")
	}
	baseUrl, err := url.Parse(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid base url: %v error: %v", baseUrl, err)
	}
	hrefUrl, err := url.Parse(selection.AttrOr("href", ""))
	if err != nil {
		return "", err
	}
	return baseUrl.ResolveReference(hrefUrl), nil
}

// Attr attr(name, defaultValue='') get element attribute value, return string.
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

// AttrConcat attrConcat(name, text1, $value, [ text2, ... text_n ])
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

// AttrEmpty attrEmpty(name, defaultValue) get element attribute value, return string.
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

// AttrSplit attrSplit(name, sep=',', trim='true')  get attribute value and split by separator to array string, return []string.
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

// EachAttr eachAttr() get each element attribute value, return []string.
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

// EachAttrEmpty eachAttrEmpty(defaultValue) get each element attribute value, return []string.
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

// EachHtml eachHtml() get each element inner html, return []string.
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

// EachOutHtml eachOutHtml() get each element outer html, return []string.
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

// EachText eachText() get each element text, return []string.
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

// EachTextEmpty eachTextEmpty(defaultValue) get each element text, return []string.
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

// EachTextJoin eachTextJoin(sep) get each element text and join to string, return string.
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

// EqAndAttr eqAndAttr(index, name) reduces the set of matched elements to the one at the specified index, and attr() return string.
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

// EqAndHtml eqAndHtml(index) reduces the set of matched elements to the one at the specified index, and html() return string.
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

// EqAndOutHtml eqAndOutHtml(index) reduces the set of matched elements to the one at the specified index, and outHtml() return string.
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

// EqAndText eqAndText(index) reduces the set of matched elements to the one at the specified index, return string.
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

// Html html() get element inner html, return string.
//	struct {
//		Example string `pagser:".selector->html()"`
//	}
func (builtin BuiltinFunctions) Html(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return node.Html()
}

// OutHtml outerHtml() get element  outer html, return string.
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

// Size size() returns the number of elements in the Selection object, return int.
//	struct {
//		Size int `pagser:".selector->size()"`
//	}
func (builtin BuiltinFunctions) Size(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return node.Size(), nil
}

// Text text() get element  text, return string, this is default function, if not define function in struct tag.
//	struct {
//		Example string `pagser:".selector->text()"`
//	}
func (builtin BuiltinFunctions) Text(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return strings.TrimSpace(node.Text()), nil
}

// TextConcat textConcat(text1, $value, [ text2, ... text_n ])
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

// TextEmpty textEmpty(defaultValue) get element text, if empty will return defaultValue, return string.
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

// TextSplit textSplit(sep=',', trim='true') get element text and split by separator to array string, return []string.
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
