# Pagser

Pagser inspired by  <u>**pag**</u>e par<u>**ser**</u>, 
an HTML deserialization to struct based on [goquery](https://github.com/PuerkitoBio/goquery) and struct tags。
It's simple, easy, extensible, configurable parsing library from [scrago](https://github.com/foolin/scrago).

## Install

```bash
go get github.com/foolin/pagser
```

## Docs

See [Pagser](https://pkg.go.dev/github.com/foolin/pagser)


## Features

* **Simple** - use golang struct tag syntax.
* **Easy** - easy use for your spider/crawler application.
* **Extensible** - Support for extension functions.
* **Struct tag** - Grammar is simple, like \`pagser:"a->attr(href)"\`.
* **Nested Structure** - Support Nested Structure for node.
* **Configurable** - Support configuration.
* **GoQuery/Colly** - Support all [goquery](https://github.com/PuerkitoBio/goquery) framework, such as [go-colly](https://github.com/gocolly/colly).

# Usage

```go

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
		Name string `pagser:"a->attr(href)"`
		Url  string `pagser:"a"`
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

```

Run output:
```

Page data json: 
-------------
{
	"Title": "Pagser Title",
	"H1": "H1 Pagser Example",
	"Navs": [
		{
			"ID": -1,
			"Name": "/",
			"Url": "Index"
		},
		{
			"ID": 2,
			"Name": "/list/web",
			"Url": "Web page"
		},
		{
			"ID": 3,
			"Name": "/list/pc",
			"Url": "Pc Page"
		},
		{
			"ID": 4,
			"Name": "/list/mobile",
			"Url": "Mobile Page"
		}
	]
}
-------------

```


# Struct tag grammar

```
[goquery selector]->[function]
```
example:



```go

type ExamData struct {
	Herf string `pagser:".navLink li a->attr(href)"`
}

selector: .navLink li a
function: attr(href)

```

# Functions

#### Builtin functions
- html()
- text()
- attr(name)
- attrInt(name, defaultValue)
- outerHtml()
- value()

#### Extensions functions
- Markdown()
- UgcHtml()

extensions function need register:
```go
import "github.com/foolin/pagser/extensions/markdown"

p := pagser.New()

//Register Markdown
markdown.Register(p)

```

#### Write self functions

```go

type CallFunc func(node *goquery.Selection, args ...string) (out interface{}, err error)

```


# Colly
Work with colly:
```colly

// On every a element which has href attribute call callback
collector.OnHTML("body", func(e *colly.HTMLElement) {
    
})

```

# Example

# Dependences
- github.com/PuerkitoBio/goquery
- github.com/spf13/cast
- github.com/mattn/godown
- github.com/microcosm-cc/bluemonday



