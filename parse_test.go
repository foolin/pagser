package pagser

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

const rawParseHtml = `
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
				<li id=""><a href="/">Index</a></li>
				<li id="2"><a href="/list/web" title="web site">Web page</a></li>
				<li id="3"><a href="/list/pc" title="pc page">Pc Page</a></li>
				<li id="4"><a href="/list/mobile" title="mobile page">Mobile Page</a></li>
			</ul>
		</div>
	</div>
	<div class="words" show="true">A|B|C|D</div>
	
	<div class="group" id="a">
		<h2>Email</h2>
		<ul>
			<li class="item" id="1" name="email" value="pagser@foolin.github">pagser@foolin.github</li>
			<li class="item" id="2" name="email" value="pagser@foolin.github">hello@pagser.foolin</li>
		</ul>
	</div>
	<div class="group" id="b">
		<h2>Bool</h2>
		<ul>
			<li class="item" id="3" name="bool" value="true">true</li>
			<li class="item" id="4" name="bool" value="false">false</li>
		</ul>
	</div>
	<div class="group" id="c">
		<h2>Number</h2>
		<ul>
			<li class="item" id="5" name="number" value="12345">12345</li>
			<li class="item" id="6" name="number" value="67890">67890</li>
		</ul>
	</div>
	<div class="group" id="d">
		<h2>Float</h2>
		<ul>
			<li class="item" id="7" name="float" value="123.45">123.45</li>
			<li class="item" id="8" name="float" value="678.90">678.90</li>
		</ul>
	</div>
</body>
</html>
`

type ParseData struct {
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
			Name   string `pagser:"->text()"`
			Url    string `pagser:"->attr(href)"`
			AbsUrl string `pagser:"->absHref('https://thisvar.com')"`
		} `pagser:"a"`
		LinkHtml       string `pagser:"a->html()"`
		ParentFuncName string `pagser:"a->ParentFunc()"`
	} `pagser:".navlink li"`
	NavFirst struct {
		ID   int    `pagser:"->attrEmpty(id, 0)"`
		Name string `pagser:"a->text()"`
		Url  string `pagser:"a->attr(href)"`
	} `pagser:".navlink li->first()"`
	NavLast struct {
		ID   int    `pagser:"->attrEmpty(id, 0)"`
		Name string `pagser:"a->text()"`
		Url  string `pagser:"a->attr(href)"`
	} `pagser:".navlink li->last()"`
	SubStruct struct {
		Label  string   `pagser:"label"`
		Values []string `pagser:".item->eachAttr(value)"`
	} `pagser:".group->eq(0)"`
	SubPtrStruct *struct {
		Label  string   `pagser:"label"`
		Values []string `pagser:".item->eachAttr(value)"`
	} `pagser:".group:last-child"`
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
	NavJoinString          string         `pagser:".navlink li->eachTextJoin(|)"`
	NavEqText              string         `pagser:".navlink li->eqAndText(1)"`
	NavEqAttr              string         `pagser:".navlink li->eqAndAttr(1, id)"`
	NavEqHtml              string         `pagser:".navlink li->eqAndHtml(1)"`
	NavEqOutHtml           string         `pagser:".navlink li->eqAndOutHtml(1)"`
	NavSize                int            `pagser:".navlink li->size()"`
	SubPageData            *SubPageData   `pagser:".navlink li:last-child"`
	SubPageDataList        []*SubPageData `pagser:".navlink li"`
	WordsSplitArray        []string       `pagser:".words->textSplit(|)"`
	WordsSplitArrayNoTrim  []string       `pagser:".words->textSplit('|', false)"`
	WordsShow              bool           `pagser:".words->attrEmpty(show, false)"`
	WordsConcatText        string         `pagser:".words->textConcat('this is words:', [, $value, ])"`
	WordsConcatAttr        string         `pagser:".words->attrConcat(show, 'isShow = [', $value, ])"`
	Email                  string         `pagser:".item[name='email']->attr('value')"`
	Emails                 []string       `pagser:".item[name='email']->eachAttrEmpty(value, '')"`
	CastBoolValue          bool           `pagser:".item[name='bool']->attrEmpty(value, false)"`
	CastBoolNoExist        bool           `pagser:".item[name='bool']->attrEmpty(value2, false)"`
	CastBoolArray          []bool         `pagser:".item[name='bool']->eachAttrEmpty(value, false)"`
	CastIntValue           int            `pagser:".item[name='number']->attrEmpty(value, 0)"`
	CastIntNoExist         int            `pagser:".item[name='number']->attrEmpty(value2, -1)"`
	CastIntArray           []int          `pagser:".item[name='number']->eachAttrEmpty(value, 0)"`
	CastInt32Value         int32          `pagser:".item[name='number']->attrEmpty(value, 0)"`
	CastInt32NoExist       int32          `pagser:".item[name='number']->attrEmpty(value2, -1)"`
	CastInt32Array         []int32        `pagser:".item[name='number']->eachAttrEmpty(value, 0)"`
	CastInt64Value         int64          `pagser:".item[name='number']->attrEmpty(value, 0)"`
	CastInt64NoExist       int64          `pagser:".item[name='number']->attrEmpty(value2, -1)"`
	CastInt64Array         []int64        `pagser:".item[name='number']->eachAttrEmpty(value, 0)"`
	CastFloat32Value       float32        `pagser:".item[name='float']->attrEmpty(value, 0)"`
	CastFloat32NoExist     float32        `pagser:".item[name='float']->attrEmpty(value2, 0.0)"`
	CastFloat32Array       []float32      `pagser:".item[name='float']->eachAttrEmpty(value, 0)"`
	CastFloat64Value       float64        `pagser:".item[name='float']->attrEmpty(value, 0)"`
	CastFloat64NoExist     float64        `pagser:".item[name='float']->attrEmpty(value2, 0.0)"`
	CastFloat64Array       []float64      `pagser:".item[name='float']->eachAttrEmpty(value, 0)"`
	NodeChild              []struct {
		Value string `pagser:"->text()"`
	} `pagser:".group->child()"`
	NodeChildSelector []struct {
		Value string `pagser:"->text()"`
	} `pagser:".group->child('h2')"`
	NodeEqFirst struct {
		Value string `pagser:"h2->text()"`
	} `pagser:".group->eq(0)"`
	NodeEqLast struct {
		Value string `pagser:"h2->text()"`
	} `pagser:".group->eq(-1)"`
	NodeEqPrev []struct {
		Value string `pagser:"->text()"`
	} `pagser:".item:last-child->prev()"`
	NodeEqPrevSelector struct {
		Value string `pagser:"->text()"`
	} `pagser:".item:last-child->prev('[id=\"1\"]')"`
	NodeEqNext []struct {
		Value string `pagser:"->text()"`
	} `pagser:".item:first-child->next()"`
	NodeEqNextSelector struct {
		Value string `pagser:"->text()"`
	} `pagser:".item:first-child->next('[id=\"2\"]')"`
	NodeParent []struct {
		Value string `pagser:"h2->text()"`
	} `pagser:"h2:first-child->parent()"`
	NodeParents []struct {
		Value string `pagser:"h2->text()"`
	} `pagser:"h2:first-child->parents()"`
	NodeParentsSelector []struct {
		Value string `pagser:"h2->text()"`
	} `pagser:"h2:first-child->parents('[id=\"b\"]')"`
	NodeParentsUntil []struct {
		Value string `pagser:"h2->text()"`
	} `pagser:"h2:first-child->parentsUntil('[id=\"b\"]')"`
	NodeParentSelector []struct {
		Value string `pagser:"h2->text()"`
	} `pagser:"h2:first-child->parent('[id=\"a\"]')"`
	NodeEqSiblings []struct {
		Value string `pagser:"->text()"`
	} `pagser:".item:first-child->siblings()"`
	NodeEqSiblingsSelector []struct {
		Value string `pagser:"->text()"`
	} `pagser:".item:first-child->siblings('[id=\"2\"]')"`
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
func (pd ParseData) MyStructFunc(selection *goquery.Selection, args ...string) (out interface{}, err error) {
	return "Struct-" + selection.Text(), nil
}

func (pd ParseData) ParentFunc(selection *goquery.Selection, args ...string) (out interface{}, err error) {
	return "ParentFunc-" + selection.Text(), nil
}

// this method will auto call, not need register.
func (pd *ParseData) FillFieldFunc(selection *goquery.Selection, args ...string) (out interface{}, err error) {
	text := selection.Text()
	pd.FillFieldOtherValue = "This value is set by the FillFieldFunc() function -" + text
	return "FillFieldFunc-" + text, nil
}

func (pd *ParseData) SameFunc(selection *goquery.Selection, args ...string) (out interface{}, err error) {
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

// Page parse from https://github.com/trending
type GithubData struct {
	Title    string `pagser:"title"`
	RepoList []struct {
		Name        string `pagser:"h1"`
		Description string `pagser:"h1 + p"`
		Stars       string `pagser:"a.muted-link->eqAndText(0)"`
		Repo        string `pagser:"h1 a->attrConcat('href', 'https://github.com', $value, '?from=pagser')"`
	} `pagser:"article.Box-row"`
}

func TestParse(t *testing.T) {
	p := New()
	//register global function
	p.RegisterFunc("MyGlobFunc", MyGlobalFunc)
	p.RegisterFunc("SameFunc", SameFunc)

	var data ParseData
	err := p.Parse(&data, rawParseHtml)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("json: %v\n", prettyJson(data))
}

func TestPagser_ParseDocument(t *testing.T) {
	cfg := Config{
		TagName:    "pagser",
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

	var data ParseData
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawParseHtml))
	if err != nil {
		t.Fatal(err)
	}
	err = p.ParseDocument(&data, doc)
	//err = p.ParseSelection(&data, doc.Selection)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("json: %v\n", prettyJson(data))
}

func TestPagser_ParseReader(t *testing.T) {
	res, err := http.Get("https://github.com/trending")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	p := New()

	var data GithubData
	err = p.ParseReader(&data, res.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("json: %v\n", prettyJson(data))
}

func TestPagser_RegisterFunc(t *testing.T) {
	threads := 1000
	p := New()
	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 10; j++ {
				for k, v := range builtinFuncs {
					p.RegisterFunc(k, v)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
