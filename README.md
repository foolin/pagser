# Pagser

[![go-doc-img]][go-doc] [![travis-img]][travis] [![go-report-card-img]][go-report-card] [![Coverage Status][cov-img]][cov]

**Pagser** inspired by  <u>**pag**</u>e par<u>**ser**</u>。

**Pagser** is a simple, easy, extensible, configurable HTML parser to struct based on [goquery](https://github.com/PuerkitoBio/goquery) and struct tags, It's parser library from [scrago](https://github.com/foolin/scrago).

## Features

* **Simple** - Use golang struct tag syntax.
* **Easy** - Easy use for your spider/crawler/colly application.
* **Extensible** - Support for extension functions.
* **Struct tag grammar** - Grammar is simple, like \`pagser:"a->attr(href)"\`.
* **Nested Structure** - Support Nested Structure for node.
* **Configurable** - Support configuration.
* **GoQuery/Colly** - Support all [goquery](https://github.com/PuerkitoBio/goquery) project, such as [go-colly](https://github.com/gocolly/colly).

## Install

```bash
go get -u github.com/foolin/pagser
```

## Docs

See [Pagser](https://pkg.go.dev/github.com/foolin/pagser)


# Usage

```go

package main

import (
	"encoding/json"
	"github.com/foolin/pagser"
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
}

func main() {
	//New default config
	p := pagser.New()

	//data parser model
	var data PageData
	//parse html data
	err := p.Parse(&data, rawPageHtml)
	//check error
	if err != nil {
		log.Fatal(err)
	}

	//print data
	log.Printf("Page data json: \n-------------\n%v\n-------------\n", toJson(data))
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
			"Name": "Index",
			"Url": "/"
		},
		{
			"ID": 2,
			"Name": "Web page",
			"Url": "/list/web"
		},
		{
			"ID": 3,
			"Name": "Pc Page",
			"Url": "/list/pc"
		},
		{
			"ID": 4,
			"Name": "Mobile Page",
			"Url": "/list/mobile"
		}
	]
}
-------------

```

# Configuration

```go

type Config struct {
	TagerName    string //struct tag name, default is `pagser`
	FuncSymbol   string //Function symbol, default is `->`
	IgnoreSymbol string //Ignore symbol, default is `-`
	Debug        bool   //Debug mode, debug will print some log, default is `false`
}

```



# Struct tag grammar

```
[goquery selector]->[function]
```
Example:
```go

type ExamData struct {
	Herf string `pagser:".navLink li a->attr(href)"`
}
```

> 1.Struct tag name: `pagser`  
> 2.[goquery](https://github.com/PuerkitoBio/goquery) selector: `.navLink li a`   
> 3.Function symbol: `->`  
> 4.Function name: `attr`  
> 5.Function arguments: `href` 

![grammar](grammar.png)

# Functions

#### Builtin functions

> - text() get element  text, return string, this is default function, if not define function in struct tag.

> - eachText() get each element text, return []string.

> - html() get element inner html, return string.

> - eachHtml() get each element inner html, return []string.

> - outerHtml() get element  outer html, return string.

> - eachOutHtml() get each element outer html, return []string.

> - attr(name) get element attribute value, return string.

> - eachAttr() get each element attribute value, return []string.

> - attrInt(name, defaultValue) get element attribute value and to int, return int.

> - attrSplit(name, sep)  get attribute value and split by separator to array string.

> - value() get element attribute value by name is `value`, return string, eg: <input value='xxxx' /> will return "xxx".

> - split(sep) get element text and split by separator to array string, return []string.

> - eachJoin(sep) get each element text and join to string, return string.


#### Extensions functions

>- Markdown() //convert html to markdown format.

>- UgcHtml() //sanitize html

Extensions function need register, like:
```go
import "github.com/foolin/pagser/extensions/markdown"

p := pagser.New()

//Register Markdown
markdown.Register(p)

```

#### Write my function

**Function interface**
```go

type CallFunc func(node *goquery.Selection, args ...string) (out interface{}, err error)

```

**1. Write global function:**
```go

//global function need call pagser.RegisterFunc("MyGlob", MyGlobalFunc) before use it.
// this global method must call pagser.RegisterFunc("MyGlob", MyGlobalFunc).
func MyGlobalFunc(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return "Global-" + node.Text(), nil
}

type PageData struct{
  MyGlobalValue string    `pagser:"->MyGlob()"`
}

func main(){

    p := pagser.New()

    //Register global function `MyGlob`
    p.RegisterFunc("MyGlob", MyGlobalFunc)

    //Todo

    //data parser model
    var data PageData
    //parse html data
    err := p.Parse(&data, rawPageHtml)

    //...
}

```


**2. Write struct function:**
```go

type PageData struct{
  MyFuncValue int    `pagser:"->MyFunc()"`
}

// this method will auto call, not need register.
func (d PageData) MyFunc(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return "Struct-" + node.Text(), nil
}


func main(){

    p := pagser.New()

    //Todo

    //data parser model
    var data PageData
    //parse html data
    err := p.Parse(&data, rawPageHtml)

    //...
}

```

**Lookup function priority order**

> struct method -> parent method -> ... -> global

See advance example: <https://github.com/foolin/pagser/tree/master/_examples/advance>

# Colly Example

Work with colly:
```go

p := pagser.New()


// On every a element which has href attribute call callback
collector.OnHTML("body", func(e *colly.HTMLElement) {
	//data parser model
	var data PageData
	//parse html data
	err := p.ParseSelection(&data, e.Dom)

})

```

# Examples

- [See Examples](https://github.com/foolin/pagser/tree/master/_examples)
- [See Tests](https://github.com/foolin/pagser/blob/master/pagser_test.go)

# Dependences
- github.com/PuerkitoBio/goquery
- github.com/spf13/cast

**Extentions:**
- github.com/mattn/godown
- github.com/microcosm-cc/bluemonday



[go-doc]: https://pkg.go.dev/github.com/foolin/pagser
[go-doc-img]: https://godoc.org/github.com/foolin/pagser?status.svg
[travis]: https://travis-ci.org/foolin/pagser
[travis-img]: https://travis-ci.org/foolin/pagser.svg?branch=master
[go-report-card]: https://goreportcard.com/report/github.com/foolin/pagser
[go-report-card-img]: https://goreportcard.com/badge/github.com/foolin/pagser
[cov-img]: https://codecov.io/gh/foolin/pagser/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/foolin/pagser