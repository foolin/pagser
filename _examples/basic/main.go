package main

import (
	"encoding/json"
	"fmt"
	"github.com/foolin/pagser"
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

type ExampPage struct {
	Title string `pagser:"title"`
	H1    string `pagser:"h1"`
	Navs  []struct {
		ID   int    `pagser:"->attrInt(id, -1)"`
		Name string `pagser:"a"`
		Url  string `pagser:"a->attrx(href)"`
	} `pagser:".navlink li"`
}

func main() {
	//New default config
	p := pagser.New()

	//data parser model
	var page ExampPage

	//parse html data
	err := p.Parse(&page, rawPageHtml)

	//check error
	if err != nil {
		panic(err)
	}

	//print data
	fmt.Printf("Page data json: \n-------------\n%v\n-------------\n", toJson(page))
}

func toJson(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}
