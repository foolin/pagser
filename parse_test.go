package pagser

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	Title               string   `pagser:"title"`
	Keywords            []string `pagser:"meta[name='keywords']->attrSplit(content)"`
	H1                  string   `pagser:"h1"`
	H1Text              string   `pagser:"h1->text()"`
	H1TextEmpty         string   `pagser:"h1->textEmpty('')"`
	H1Html              string   `pagser:"h1->html()"`
	H1OutHtml           string   `pagser:"h1->outerHtml()"`
	SameFuncValue       string   `pagser:"h1->SameFunc()"`
	MyGlobalFuncValue   string   `pagser:"h1->MyGlobFunc()"`
	MyStructFuncValue   string   `pagser:"h1->MyStructFunc()"`
	FillFieldFuncValue  string   `pagser:"h1->FillFieldFunc()"`
	FillFieldOtherValue string   //Set value by FillFieldFunc()
	NavList             []struct {
		ID   int `pagser:"->attrEmpty(id, -1)"`
		Link struct {
			Name string `pagser:"->text()"`
			Url  string `pagser:"->attr(href)"`
		} `pagser:"a"`
		LinkHtml       string `pagser:"a->html()"`
		ParentFuncName string `pagser:"a->ParentFunc()"`
	} `pagser:".navlink li"`
	NavLast struct {
		ID   int    `pagser:"->attrEmpty(id, 0)"`
		Name string `pagser:"a->text()"`
		Url  string `pagser:"a->attr(href)"`
	} `pagser:".navlink li:last-child"`
	NavFirstID             int            `pagser:".navlink li:first-child->attrEmpty(id, 0)"`
	NavLastID              uint           `pagser:".navlink li:last-child->attr(id)"`
	NavFirstIDDefaultValue int            `pagser:".navlink li:first-child->attrEmpty(id, -999)"`
	NavTextList            []string       `pagser:".navlink li"`
	NavEachText            []string       `pagser:".navlink li->eachText()"`
	NavEachTextEmpty       []string       `pagser:".navlink li->eachTextEmpty('')"`
	NavEachAttrID          []string       `pagser:".navlink li->eachAttr(id)"`
	NavEachAttrEmptyID     []string       `pagser:".navlink li->eachAttrEmpty(id, -1)"`
	NavEachHtml            []string       `pagser:".navlink li->eachHtml()"`
	NavEachOutHtml         []string       `pagser:".navlink li->eachOutHtml()"`
	NavJoinString          string         `pagser:".navlink li->eachJoin(|)"`
	NavEq                  string         `pagser:".navlink li->eq(1)"`
	NavEqAttr              string         `pagser:".navlink li->eqAndAttr(1, id)"`
	NavEqHtml              string         `pagser:".navlink li->eqAndHtml(1)"`
	NavEqOutHtml           string         `pagser:".navlink li->eqAndOutHtml(1)"`
	SubPageData            *SubPageData   `pagser:".navlink li:last-child"`
	SubPageDataList        []*SubPageData `pagser:".navlink li"`
	WordsSplitArray        []string       `pagser:".words->split(|)"`
	WordsShow              bool           `pagser:".words->attrEmpty(show, false)"`
	WordsConcat            string         `pagser:".words->concat('this is words:', [, $value, ])"`
	WordsConcatAttr        string         `pagser:".words->concatAttr(show, 'isShow = [', $value, ])"`
	Email                  string         `pagser:"input[name='email']->value()"`
	Emails                 []string       `pagser:"input[name='email']->eachAttrEmpty(value, '')"`
	CastBoolValue          bool           `pagser:"input[name='bool']->attrEmpty(value, false)"`
	CastBoolNoExist        bool           `pagser:"input[name='bool']->attrEmpty(value2, false)"`
	CastBoolArray          []bool         `pagser:"input[name='bool']->eachAttrEmpty(value, false)"`
	CastIntValue           int            `pagser:"input[name='number']->attrEmpty(value, 0)"`
	CastIntNoExist         int            `pagser:"input[name='number']->attrEmpty(value2, -1)"`
	CastIntArray           []int          `pagser:"input[name='number']->eachAttrEmpty(value, 0)"`
	CastInt32Value         int32          `pagser:"input[name='number']->attrEmpty(value, 0)"`
	CastInt32NoExist       int32          `pagser:"input[name='number']->attrEmpty(value2, -1)"`
	CastInt32Array         []int32        `pagser:"input[name='number']->eachAttrEmpty(value, 0)"`
	CastInt64Value         int64          `pagser:"input[name='number']->attrEmpty(value, 0)"`
	CastInt64NoExist       int64          `pagser:"input[name='number']->attrEmpty(value2, -1)"`
	CastInt64Array         []int64        `pagser:"input[name='number']->eachAttrEmpty(value, 0)"`
	CastFloat32Value       float32        `pagser:"input[name='float']->attrEmpty(value, 0)"`
	CastFloat32NoExist     float32        `pagser:"input[name='float']->attrEmpty(value2, 0.0)"`
	CastFloat32Array       []float32      `pagser:"input[name='float']->eachAttrEmpty(value, 0)"`
	CastFloat64Value       float64        `pagser:"input[name='float']->attrEmpty(value, 0)"`
	CastFloat64NoExist     float64        `pagser:"input[name='float']->attrEmpty(value2, 0.0)"`
	CastFloat64Array       []float64      `pagser:"input[name='float']->eachAttrEmpty(value, 0)"`
}

// this method will auto call, not need register.
func MyGlobalFunc(selection *goquery.Selection, args ...string) (out interface{}, err error) {
	return "Global-" + selection.Text(), nil
}

// this method will auto call, not need register.
func SameFunc(selection *goquery.Selection, args ...string) (out interface{}, err error) {
	return "Global-Same-Func-" + selection.Text(), nil
}

// this method will auto call, not need register.
func (pd PageData) MyStructFunc(selection *goquery.Selection, args ...string) (out interface{}, err error) {
	return "Struct-" + selection.Text(), nil
}

func (pd PageData) ParentFunc(selection *goquery.Selection, args ...string) (out interface{}, err error) {
	return "ParentFunc-" + selection.Text(), nil
}

// this method will auto call, not need register.
func (pd *PageData) FillFieldFunc(selection *goquery.Selection, args ...string) (out interface{}, err error) {
	text := selection.Text()
	pd.FillFieldOtherValue = "This value is set by the FillFieldFunc() function -" + text
	return "FillFieldFunc-" + text, nil
}

func (pd *PageData) SameFunc(selection *goquery.Selection, args ...string) (out interface{}, err error) {
	return "Struct-Same-Func-" + selection.Text(), nil
}

type SubPageData struct {
	Text            string `pagser:"->text()"`
	SubFuncValue    string `pagser:"->SubFunc()"`
	ParentFuncValue string `pagser:"->ParentFunc()"`
	SameFuncValue   string `pagser:"->SameFunc()"`
}

func (spd SubPageData) SubFunc(selection *goquery.Selection, args ...string) (out interface{}, err error) {
	return "SubFunc-" + selection.Text(), nil
}

func (spd SubPageData) SameFunc(selection *goquery.Selection, args ...string) (out interface{}, err error) {
	return "Sub-Struct-Same-Func-" + selection.Text(), nil
}

const rawPpageHtml = `
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>Pagser Example</title>
	<meta name="keywords" content="golang,pagser,goquery,html,page,parser,colly">
</head>

<body>
	<h1><u>Pagser</u> H1 Title</h1>
	<div class="navlink">
		<div class="container">
			<ul class="clearfix">
				<li id=''><a href="/">Index</a></li>
				<li id='2'><a href="/list/web" title="web site">Web page</a></li>
				<li id='3'><a href="/list/pc" title="pc page">Pc Page</a></li>
				<li id='4'><a href="/list/mobile" title="mobile page">Mobile Page</a></li>
			</ul>
		</div>
	</div>
	<div class='words' show="true">A|B|C|D</div>
	
	<input name="email" value="pagser@foolin.github" />
	<input name="email" value="hello@pagser.foolin" />
	<input name="bool" value="true" /> 
	<input name="bool" value="false" /> 
	<input name="number" value="12345" />
	<input name="number" value="67890" />
	<input name="float" value="123.45" /> 
	<input name="float" value="678.90" />
</body>
</html>
`

func TestParse(t *testing.T) {
	p := New()
	//register global function
	p.RegisterFunc("MyGlobFunc", MyGlobalFunc)
	p.RegisterFunc("SameFunc", SameFunc)

	var data PageData
	err := p.Parse(&data, rawPpageHtml)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("json: %v\n", prettyJson(data))
}

func TestPagser_ParseDocument(t *testing.T) {
	cfg := Config{
		TagerName:  "pagser",
		FuncSymbol: "->",
		CastError:  true,
		Debug:      true,
	}
	p, err := NewWithConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}

	//register global function
	p.RegisterFunc("MyGlobFunc", MyGlobalFunc)
	p.RegisterFunc("SameFunc", SameFunc)

	var data PageData
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawPageHtml))
	if err != nil {
		t.Fatal(err)
	}
	err = p.ParseDocument(&data, doc)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("json: %v\n", prettyJson(data))
}

func TestPagser_ParseReader(t *testing.T) {
	res, err := http.Get("https://raw.githubusercontent.com/foolin/pagser/master/_examples/pages/demo.html")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	p := New()

	//register global function
	p.RegisterFunc("MyGlobFunc", MyGlobalFunc)

	var data PageData
	err = p.ParseReader(&data, res.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("json: %v\n", prettyJson(data))
}
