package pagser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"testing"
)

//	"html":      html,
//	"text":      text,
//	"attr":      attr,
//	"attrInt":   attrInt,
//	"outerHtml": outHtml,
//	"value":     value,
type HtmlPage struct {
	Title             string   `pagser:"title"`
	Keywords          []string `pagser:"meta[name='keywords']->attrSplit(content)"`
	H1                string   `pagser:"h1->text()"`
	H1Html            string   `pagser:"h1->html()"`
	H1OutHtml         string   `pagser:"h1->outerHtml()"`
	MyGlobalFuncValue string   `pagser:"h1->MyGlobFunc()"`
	MyStructFuncValue string   `pagser:"h1->MyStructFunc()"`
	NavList           []struct {
		ID   int    `pagser:"a->attrInt(id, -1)"`
		Name string `pagser:"a"`
		Url  string `pagser:"a->attr(href)"`
	} `pagser:".navlink li"`
	NavTextList     []string `pagser:".navlink li"`
	NavEachText     []string `pagser:".navlink li->eachText()"`
	NavEachAttrID   []string `pagser:".navlink li->eachAttr(id)"`
	NavEachHtml     []string `pagser:".navlink li->eachHtml()"`
	NavEachOutHtml  []string `pagser:".navlink li->eachOutHtml()"`
	NavJoinString   string   `pagser:".navlink li->eachJoin(|)"`
	WordsSplitArray []string `pagser:".words->split(|)"`
	Value           string   `pagser:"input[name='feedback']->value()"`
}

// this method will auto call, not need register.
func MyGlobalFunc(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return "Global-" + node.Text(), nil
}

// this method will auto call, not need register.
func (h HtmlPage) MyStructFunc(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return "Struct-" + node.Text(), nil
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
	<div class='words'>A|B|C|D</div>
	<input name="feedback" value="pagser@foolin.github" /> 
</body>
</html>
`

func TestParse(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Debug = true
	p, err := NewWithConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}

	//register global function
	p.RegisterFunc("MyGlobFunc", MyGlobalFunc)

	var data HtmlPage
	err = p.Parse(&data, rawPpageHtml)
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

	var data HtmlPage
	err = p.ParseReader(&data, res.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("json: %v\n", prettyJson(data))
}

func TestParseFuncArgs(t *testing.T) {
	//const str = `"Hello world, Andy", abc, 123, 1.2,3, -1, ',','abc\'x', -12.3, true, '1,2,3', false, "[1,2,3]",bool`
	//const str = `Foo bar random "letters lol" stuff`
	const str = `href,123, 456, 'Foolin said: "it\'s ok"', "not work", 'Hello world'`
	//const str = "href"
	params := parseFuncParams(str)
	for i, v := range params {
		fmt.Printf("%.2d: %v\n", i, v)
	}
	fmt.Printf("%v", prettyJson(params))
}
