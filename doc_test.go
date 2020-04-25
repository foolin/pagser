package pagser

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

const rawExampleHtml = `
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

type ExamplePage struct {
	Title string `pagser:"title"`
	H1    string `pagser:"h1"`
	Navs  []struct {
		ID   int    `pagser:"->attrEmpty(id, -1)"`
		Name string `pagser:"a"`
		Url  string `pagser:"a->attr(href)"`
	} `pagser:".navlink li"`
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
	var page ExamplePage
	//parse html data
	err = p.Parse(&page, rawExampleHtml)
	//check error
	if err != nil {
		log.Fatal(err)
	}

}

func ExamplePagser_Parse() {
	//New default Config
	p := New()

	//data parser model
	var page ExamplePage
	//parse html data
	err := p.Parse(&page, rawExampleHtml)
	//check error
	if err != nil {
		log.Fatal(err)
	}

	//print result
	log.Printf("%v", page)
}

func ExamplePagser_ParseDocument() {
	//New default Config
	p := New()

	//data parser model
	var data ExamplePage
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawExampleHtml))
	if err != nil {
		log.Fatal(err)
	}

	//parse document
	err = p.ParseDocument(&data, doc)
	//check error
	if err != nil {
		log.Fatal(err)
	}

	//print result
	log.Printf("%v", data)
}

func ExamplePagser_ParseSelection() {
	//New default Config
	p := New()

	//data parser model
	var data ExamplePage
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawExampleHtml))
	if err != nil {
		log.Fatal(err)
	}

	//parse document
	err = p.ParseSelection(&data, doc.Selection)
	//check error
	if err != nil {
		log.Fatal(err)
	}

	//print result
	log.Printf("%v", data)
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
	var page ExamplePage
	//parse html data
	err = p.ParseReader(&page, resp.Body)
	//check error
	if err != nil {
		panic(err)
	}

	log.Printf("%v", page)
}

func ExamplePagser_RegisterFunc() {
	p := New()

	p.RegisterFunc("MyFunc", func(node *goquery.Selection, args ...string) (out interface{}, err error) {
		//Todo
		return "Hello", nil
	})
}
