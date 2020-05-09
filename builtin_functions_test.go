package pagser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"testing"
)

type funcWantError struct {
	want bool
	fun  string
	args []string
	data string
}

func (fwe funcWantError) String() string {
	return fmt.Sprintf("%v(%v) call `%v`", fwe.fun, strings.Join(fwe.args, ","), fwe.data)
}

func newTewSelection(content string) *goquery.Selection {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(content))
	return doc.Selection.Find("body").Children()
}

func TestBuiltinFunctionsErrors(t *testing.T) {
	tests := []funcWantError{
		//not baseUrl
		{true, "absHref", []string{}, ``},
		//baseUrl invalid
		{true, "absHref", []string{"http://a b.com/"}, `<a href="/foo/bar">a</a>`},
		//href value not invalid url
		{true, "absHref", []string{"http://github.com/"}, `<a href="cache_object:foo/bar">a</a>`},
		//not attr name
		{true, "attr", []string{}, `<a href="/foo">a</a>`},
		//not attr name
		{true, "attrConcat", []string{}, `<a href="/foo">a</a>`},
		//not attr $value
		{true, "attrConcat", []string{"href"}, `<a href="/foo">a</a>`},
		//not concat second string
		{true, "attrConcat", []string{"href", "$value"}, `<a href="/foo">a</a>`},
		//not default value
		{true, "attrEmpty", []string{"href"}, `<a href="/foo">a</a>`},
		//not attr name
		{true, "attrSplit", []string{}, `<a href="/foo">a</a>`},
		//trim value '1234' is not bool type
		{true, "attrSplit", []string{"href", "|", "1234"}, `<a href="/foo">a</a>`},
		//not attr name
		{true, "eachAttr", []string{}, `<a href="/foo">a</a>`},
		//not default value
		{true, "eachAttrEmpty", []string{"href"}, `<a href="/foo">a</a>`},
		//html error
		//{true, "eachOutHtml", []string{}, `</a>abc<a>`},
		//not default value
		{true, "eachTextEmpty", []string{}, `<a href="/foo">a</a>`},
		//not index value
		{true, "eqAndAttr", []string{}, `<a href="/foo">a</a>`},
		//not name value
		{true, "eqAndAttr", []string{"0"}, `<a href="/foo">a</a>`},
		//index not number
		{true, "eqAndAttr", []string{"a", "href"}, `<a href="/foo">a</a>`},
		//not index value
		{true, "eqAndHtml", []string{}, `<a href="/foo">a</a>`},
		//index not number
		{true, "eqAndHtml", []string{"a"}, `<a href="/foo">a</a>`},
		//not index value
		{true, "eqAndOutHtml", []string{}, `<a href="/foo">a</a>`},
		//not name value
		{true, "eqAndOutHtml", []string{"a"}, `<a href="/foo">a</a>`},
		//not index value
		{true, "eqAndText", []string{}, `<a href="/foo">a</a>`},
		//not name value
		{true, "eqAndText", []string{"a"}, `<a href="/foo">a</a>`},
		//html error
		//{true, "outerHtml", []string{"a"}, `</aa</a<>`},
		//not args
		{true, "textConcat", []string{"$value"}, `</a>a</a>`},
		//not args
		{true, "textEmpty", []string{}, `<a href="/foo">a</a>`},
		//not bool
		{true, "textSplit", []string{",", "1.2"}, `<a href="/foo">a</a>`},
	}

	for _, tt := range tests {
		var sel = newTewSelection(tt.data)
		_, err := builtinFuncs[tt.fun](sel, tt.args...)
		if tt.want {
			if err == nil {
				t.Errorf("%v want an error", tt.String())
			}
			continue
		}
		if err != nil {
			t.Errorf("%v want no error, but error is %v", tt.String(), err)
		}
	}
}
