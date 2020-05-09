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
		{true, "absHref", []string{}, ``},
		{true, "absHref", []string{"http://a b.com/"}, `<a href="/hello/foolin">a</a>`},
		{true, "absHref", []string{"http://github.com/"}, `<a href="cache_object:foo/bar">a</a>`},
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
