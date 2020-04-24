package pagser

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

const rawPageHtml = `
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>Pagser Title</title>
	<meta name="keywords" content="golang,pagser,goquery,html,page,parser,colly">
</head>

<body>
	<h1>H1 Pagser Example</h1>
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
</body>
</html>
`

type ExampPage struct {
	Title string `pagser:"title"`
	H1    string `pagser:"h1"`
	Navs  []struct {
		ID   int    `pagser:"->attrEmpty(id, -1)"`
		Name string `pagser:"a"`
		Url  string `pagser:"a->attr(href)"`
	} `pagser:".navlink li"`
}

func ExamplePagser_Parse() {
	//New default Config
	p := New()
	//data parser model
	var page ExampPage
	//parse html data
	err := p.Parse(&page, rawPageHtml)
	//check error
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v", page)
}

func ExamplePagser_ParseReader() {
	resp, err := http.Get("https://raw.githubusercontent.com/foolin/pagser/master/_examples/pages/demo.html")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//New default Config
	p := New()
	//data parser model
	var page ExampPage
	//parse html data
	err = p.ParseReader(&page, resp.Body)
	//check error
	if err != nil {
		panic(err)
	}

	log.Printf("%v", page)
}

func ExampleNewWithConfig() {
	cfg := Config{
		TagName:    "pagser",
		FuncSymbol: "->",
		CastError:  false,
		Debug:      false,
	}
	p, err := NewWithConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	//data parser model
	var page ExampPage
	//parse html data
	err = p.Parse(&page, rawPageHtml)
	//check error
	if err != nil {
		log.Fatal(err)
	}

}

func MyFunc(node *goquery.Selection, args ...string) (out interface{}, err error) {
	//todo
	return "Hello", nil
}
func ExamplePagser_RegisterFunc() {
	p := New()
	p.RegisterFunc("MyFunc", MyFunc)
}
