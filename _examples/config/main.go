package main

import (
	"encoding/json"
	"github.com/foolin/pagser"
	"log"
)

type ExampleData struct {
	Title    string   `query:"title"`
	Keywords []string `query:"meta[name='keywords']@attrSplit(content)"`
	Navs     []struct {
		ID   int    `query:"@attrEmpty(id, -1)"`
		Name string `query:"a@text()"`
		Url  string `query:"a@attr(href)"`
	} `query:".navlink li"`
}

const rawExampleHtml = `
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>Pagser Example</title>
	<meta name="keywords" content="golang,pagser,goquery,html,page,parser,colly">
</head>

<body>
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
	cfg := pagser.Config{
		TagerName:  "query",
		FuncSymbol: "@",
		CastError:  true,
		Debug:      true,
	}
	p, err := pagser.NewWithConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	var data ExampleData
	err = p.Parse(&data, rawExampleHtml)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("json: %v\n", toJson(data))
}

func toJson(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}
