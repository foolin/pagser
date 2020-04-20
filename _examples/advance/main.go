package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/foolin/pagser"
	"github.com/foolin/pagser/extensions/markdown"
	"strconv"
	"strings"
)

const rawPageHtml = `
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>PagserTitle</title>
</head>

<body>
	<h1>Test H1 Example</h1>
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

func main() {
	//New with config
	cfg := pagser.DefaultConfig()
	cfg.Debug = true
	p := pagser.MustNewWithConfig(cfg)

	//Register Markdown
	markdown.Register(p)
	//Register AttrInt
	p.RegisterFunc("AttrInt", AttrInt)

	var page ExampPage
	err := p.Parse(&page, rawPageHtml)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Page data json: \n-------------\n%v\n-------------\n", toJson(page))
}

type ExampPage struct {
	Title string `pagser:"title"`
	H1    string `pagser:"h1"`
	Navs  []struct {
		ID   int    `pagser:"->attrInt(id, -1)"`
		Name string `pagser:"a"`
		Url  string `pagser:"a->attr(href)"`
	} `pagser:".navlink li"`
	NavID    []int    `pagser:".navlink li->AttrInt(id, '-1')"`
	NavTexts []string `pagser:".navlink li"`
	NavMods  []string `pagser:".navlink li->GetMod()"`
	Markdown string   `pagser:"->Markdown()"`
}

// this method will auto call, not need register.
func (h ExampPage) GetMod(node *goquery.Selection, args ...string) (out interface{}, err error) {
	//<li id='3'><a href="/list/pc" title="pc page">Pc Page</a></li>
	href := node.Find("a").AttrOr("href", "")
	fmt.Printf("href: %v\n", href)
	if idx := strings.LastIndex(href, "/"); idx != -1 {
		return href[idx+1:], nil
	}
	return "", nil
}

// this global method must call pagser.RegisterFunc("AttrInt", AttrInt).
func AttrInt(node *goquery.Selection, args ...string) (out interface{}, err error) {
	//<li id='2'>
	//`pagser:".navlink li->AttrInt(id, '-1')"`
	if len(args) < 2 {
		return -1, errors.New("AttrInt must has two args, such as AttrInt(id, 0)")
	}
	value := node.AttrOr(args[0], "")
	if value == "" {
		return -1, nil
	}
	return strconv.Atoi(value)
}

func toJson(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}
