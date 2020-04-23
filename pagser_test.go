package pagser

import (
	"fmt"
	"testing"
)

type ExampleData struct {
	Title    string   `pagser:"title"`
	Keywords []string `pagser:"meta[name='keywords']->attrSplit(content)"`
	Navs     []struct {
		ID   int    `pagser:"->attrEmpty(id, -1)"`
		Name string `pagser:"a->text()"`
		Url  string `pagser:"a->attr(href)"`
	} `pagser:".navlink li"`
}

type ConfigData struct {
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

func TestNew(t *testing.T) {
	p := New()

	var data ExampleData
	err := p.Parse(&data, rawExampleHtml)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("json: %v\n", prettyJson(data))
}

func TestNewWithConfig(t *testing.T) {
	cfg := Config{
		TagerName:  "query",
		FuncSymbol: "@",
		CastError:  true,
		Debug:      true,
	}
	p, err := NewWithConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}

	var data ConfigData
	err = p.Parse(&data, rawExampleHtml)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("json: %v\n", prettyJson(data))
}
