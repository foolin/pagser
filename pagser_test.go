package pagser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"testing"
)

//	"html":      html,
//	"text":      text,
//	"attr":      attr,
//	"attrInt":   attrInt,
//	"outerHtml": outHtml,
//	"value":     value,
type HtmlPage struct {
	Title  string `pagser:"title"`
	H1     string `pagser:"h1->text()"`
	H1HTML string `pagser:"h1->outerHtml()"`
	Navs   []struct {
		ID   int    `pagser:"a->attrInt(id, -1)"`
		Name string `pagser:"a->attr(href)"`
		Url  string `pagser:"a"`
	} `pagser:".navlink li"`
	NavNames   []string `pagser:".navlink li"`
	NavMods    []string `pagser:".navlink li->GetMod()"`
	FirstLink  string   `pagser:".navlink li:first-child->html()"`
	InputValue string   `pagser:"input[name='feedback']->value()"`
}

// this method will auto call, not need register.
func (h HtmlPage) GetMod(node *goquery.Selection, args ...string) (out interface{}, err error) {
	//<li id='3'><a href="/list/pc" title="pc page">Pc Page</a></li>
	href := node.Find("a").AttrOr("href", "")
	if idx := strings.LastIndex(href, "/"); idx != -1 {
		return href[idx+1:], nil
	}
	return "", nil
}

const rawPpageHtml = `
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>Pagser Example</title>
</head>

<body>
	<h1>Pagser H1 Title</h1>
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
	<input name="feedback" value="hello@example.com" /> 
</body>
</html>
`

func TestParse(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Debug = true
	p := MustNewWithConfig(cfg)

	var data HtmlPage
	err := p.Parse(&data, rawPpageHtml)
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
