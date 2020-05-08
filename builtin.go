package pagser

import (
	"github.com/PuerkitoBio/goquery"
)

// CallFunc write function interface
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

//BuiltinFunctions instance
var builtinFun BuiltinFunctions

//BuiltinSelections instance
var builtinSel BuiltinSelections

//builtin functions
var builtinFuncs = map[string]CallFunc{
	"absHref":       builtinFun.AbsHref,
	"attr":          builtinFun.Attr,
	"attrConcat":    builtinFun.AttrConcat,
	"attrEmpty":     builtinFun.AttrEmpty,
	"attrSplit":     builtinFun.AttrSplit,
	"eachAttr":      builtinFun.EachAttr,
	"eachAttrEmpty": builtinFun.EachAttrEmpty,
	"eachHtml":      builtinFun.EachHtml,
	"eachOutHtml":   builtinFun.EachOutHtml,
	"eachText":      builtinFun.EachText,
	"eachTextEmpty": builtinFun.EachTextEmpty,
	"eachTextJoin":  builtinFun.EachTextJoin,
	"eqAndAttr":     builtinFun.EqAndAttr,
	"eqAndHtml":     builtinFun.EqAndHtml,
	"eqAndOutHtml":  builtinFun.EqAndOutHtml,
	"eqAndText":     builtinFun.EqAndText,
	"html":          builtinFun.Html,
	"outerHtml":     builtinFun.OutHtml,
	"size":          builtinFun.Size,
	"text":          builtinFun.Text,
	"textConcat":    builtinFun.TextConcat,
	"textEmpty":     builtinFun.TextEmpty,
	"textSplit":     builtinFun.TextSplit,
	// selector
	"child":        builtinSel.Child,
	"eq":           builtinSel.Eq,
	"first":        builtinSel.First,
	"last":         builtinSel.Last,
	"next":         builtinSel.Next,
	"parent":       builtinSel.Parent,
	"parents":      builtinSel.Parents,
	"parentsUntil": builtinSel.ParentsUntil,
	"prev":         builtinSel.Prev,
	"siblings":     builtinSel.Siblings,
}

// RegisterFunc register function for parse result
//	pagser.RegisterFunc("MyFunc", func(node *goquery.Selection, args ...string) (out interface{}, err error) {
//		//Todo
//		return "Hello", nil
//	})
func (p *Pagser) RegisterFunc(name string, fn CallFunc) {
	p.mapFuncs.Store(name, fn)
}
