package main

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/foolin/pagser"
	"github.com/foolin/pagser/extensions/markdown"
	"log"
)

const rawPageHtml = `
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>Pagser Title</title>
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

type PageData struct {
	Title string `pagser:"title"`
	H1    string `pagser:"h1"`
	Navs  []struct {
		ID   int    `pagser:"->attrInt(id, -1)"`
		Name string `pagser:"a"`
		Url  string `pagser:"a->attr(href)"`
	} `pagser:".navlink li"`
	NavIds        []string `pagser:".navlink li->eachAttr(id, -1)"`
	NavTexts      []string `pagser:".navlink li"`
	NavEachTexts  []string `pagser:".navlink li->eachText()"`
	MyFuncValue   string   `pagser:"h1->MyFunc()"`
	MyGlobalValue string   `pagser:"h1->MyGlob()"`
	Markdown      string   `pagser:"->Markdown()"`
}

// this method will auto call, not need register.
func (pd PageData) MyFunc(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return "Struct-" + node.Text(), nil
}

// this global method must call pagser.RegisterFunc("MyGlobFunc", MyGlobalFunc).
func MyGlobalFunc(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return "Global-" + node.Text(), nil
}

func main() {
	//New with config
	cfg := pagser.Config{
		TagerName:    "pagser",
		FuncSymbol:   "->",
		IgnoreSymbol: "-",
		Debug:        true,
	}
	p, err := pagser.NewWithConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	//Register Markdown
	markdown.Register(p)

	//Register global function
	p.RegisterFunc("MyGlob", MyGlobalFunc)

	var data PageData
	err = p.Parse(&data, rawPageHtml)
	if err != nil {
		panic(err)
	}

	log.Printf("Page data json: \n-------------\n%v\n-------------\n", toJson(data))
}

func toJson(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}
