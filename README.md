# Pagser

[![go-doc-img]][go-doc] [![travis-img]][travis] [![go-report-card-img]][go-report-card] [![Coverage Status][cov-img]][cov]

**Pagser** inspired by  <u>**pag**</u>e par<u>**ser**</u>。

**Pagser** is a simple, easy, extensible, configurable HTML parser to struct based on [goquery](https://github.com/PuerkitoBio/goquery) and struct tags, It's parser library from [scrago](https://github.com/foolin/scrago).


## Contents

- [Install](#install)
- [Features](#features)
- [Docs](#docs)
- [Configuration](#configuration)
- [Usage](#usage)
- [Struct Tag Grammar](#struct-tag-grammar)
- [Functions](#functions)
    - [Builtin functions](#builtin-functions)
    - [Extension functions](#extension-functions)
    - [Custom function](#custom-function)
    - [Function interface](#function-interface)
    - [Call Syntax](#call-syntax)
    - [Priority Order](#priority-order)
    - [More Examples](#more-examples)
- [Examples](#examples)
- [Dependencies](#dependencies)


## Install

```bash
go get -u github.com/foolin/pagser
```

## Features

* **Simple** - Use golang struct tag syntax.
* **Easy** - Easy use for your spider/crawler/colly application.
* **Extensible** - Support for extension functions.
* **Struct tag grammar** - Grammar is simple, like \`pagser:"a->attr(href)"\`.
* **Nested Structure** - Support Nested Structure for node.
* **Configurable** - Support configuration.
* **Implicit type conversion** - Automatic implicit type conversion, Output result string convert to int, int64, float64...
* **GoQuery/Colly** - Support all [goquery](https://github.com/PuerkitoBio/goquery) project, such as [go-colly](https://github.com/gocolly/colly).

## Docs

See [Pagser](https://pkg.go.dev/github.com/foolin/pagser)


## Usage

```golang

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

type PageData struct {
	Title    string   `pagser:"title"`
	Keywords []string `pagser:"meta[name='keywords']->attrSplit(content)"`
	H1       string   `pagser:"h1"`
	Navs     []struct {
		ID   int    `pagser:"->attrEmpty(id, -1)"`
		Name string `pagser:"a->text()"`
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
	"Keywords": [
		"golang",
		"pagser",
		"goquery",
		"html",
		"page",
		"parser",
		"colly"
	],
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

## Configuration

```golang

type Config struct {
	TagName    string //struct tag name, default is `pagser`
	FuncSymbol   string //Function symbol, default is `->`
	Debug        bool   //Debug mode, debug will print some log, default is `false`
}

```



## Struct Tag Grammar

```
[goquery selector]->[function]
```
Example:
```golang

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

## Functions

### Builtin functions

> - text() get element  text, return string, this is default function, if not define function in struct tag.

> - eachText() get each element text, return []string.

> - html() get element inner html, return string.

> - eachHtml() get each element inner html, return []string.

> - outerHtml() get element  outer html, return string.

> - eachOutHtml() get each element outer html, return []string.

> - attr(name) get element attribute value, return string.

> - eachAttr() get each element attribute value, return []string.

> - attrSplit(name, sep)  get attribute value and split by separator to array string.

> - value() get element attribute value by name is `value`, return string, eg: <input value='xxxx' /> will return "xxx".

> - split(sep) get element text and split by separator to array string, return []string.

> - eachJoin(sep) get each element text and join to string, return string.

> - ...

More builtin functions see docs: <https://pkg.go.dev/github.com/foolin/pagser>

### Extension functions

>- Markdown() //convert html to markdown format.

>- UgcHtml() //sanitize html

Extensions function need register, like:
```golang
import "github.com/foolin/pagser/extensions/markdown"

p := pagser.New()

//Register Markdown
markdown.Register(p)

```

### Custom function

#### Function interface
```golang

type CallFunc func(node *goquery.Selection, args ...string) (out interface{}, err error)

```

#### Define global function
```golang

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


#### Define struct function
```golang

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

#### Call Syntax

> **Note**: all function arguments are string, single quotes are optional.

1. Function call with no arguments
> ->fn()

2. Function calls with one argument, and single quotes are optional

> ->fn(one)
>
> ->fn('one')

3. Function calls with many arguments

> ->fn(one, two, three, ...)
>
> ->fn('one', 'two', 'three', ...)


5. Function calls with single quotes and escape character

> ->fn('it\\'s ok', 'two,xxx', 'three', ...)


### Priority Order

Lookup function priority order:

> struct method -> parent method -> ... -> global


### More Examples
See advance example: <https://github.com/foolin/pagser/tree/master/_examples/advance>

## Implicit type conversion
Automatic implicit type conversion, Output result string convert to int, int64, float64...

**Support type:**

- bool
- float32
- float64
- int
- int32
- int64
- string
- []bool
- []float32
- []float64
- []int
- []int32
- []int64
- []string



## Examples

### Crawl page example

```golang

package main

import (
	"encoding/json"
	"github.com/foolin/pagser"
	"log"
	"net/http"
)

type PageData struct {
	Title    string `pagser:"title"`
	RepoList []struct {
		Names       []string `pagser:"h1->split('/', true)"`
		Description string   `pagser:"h1 + p"`
		Stars       string   `pagser:"a.muted-link->eq(0)"`
		Repo        string   `pagser:"h1 a->concatAttr('href', 'https://github.com', $value, '?from=pagser')"`
	} `pagser:"article.Box-row"`
}

func main() {
	resp, err := http.Get("https://github.com/trending")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//New default config
	p := pagser.New()

	//data parser model
	var data PageData
	//parse html data
	err = p.ParseReader(&data, resp.Body)
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

2020/04/25 12:26:04 Page data json: 
-------------
{
	"Title": "Trending  repositories on GitHub today · GitHub",
	"RepoList": [
		{
			"Names": [
				"pcottle",
				"learnGitBranching"
			],
			"Description": "An interactive git visualization to challenge and educate!",
			"Stars": "16,010",
			"Repo": "https://github.com/pcottle/learnGitBranching?from=pagser"
		},
		{
			"Names": [
				"jackfrued",
				"Python-100-Days"
			],
			"Description": "Python - 100天从新手到大师",
			"Stars": "83,484",
			"Repo": "https://github.com/jackfrued/Python-100-Days?from=pagser"
		},
		{
			"Names": [
				"brave",
				"brave-browser"
			],
			"Description": "Next generation Brave browser for macOS, Windows, Linux, Android.",
			"Stars": "5,963",
			"Repo": "https://github.com/brave/brave-browser?from=pagser"
		},
		{
			"Names": [
				"MicrosoftDocs",
				"azure-docs"
			],
			"Description": "Open source documentation of Microsoft Azure",
			"Stars": "3,798",
			"Repo": "https://github.com/MicrosoftDocs/azure-docs?from=pagser"
		},
		{
			"Names": [
				"ahmetb",
				"kubectx"
			],
			"Description": "Faster way to switch between clusters and namespaces in kubectl",
			"Stars": "6,979",
			"Repo": "https://github.com/ahmetb/kubectx?from=pagser"
		},

        //...        

		{
			"Names": [
				"serverless",
				"serverless"
			],
			"Description": "Serverless Framework – Build web, mobile and IoT applications with serverless architectures using AWS Lambda, Azure Functions, Google CloudFunctions \u0026 more! –",
			"Stars": "35,502",
			"Repo": "https://github.com/serverless/serverless?from=pagser"
		},
		{
			"Names": [
				"vuejs",
				"vite"
			],
			"Description": "Experimental no-bundle dev server for Vue SFCs",
			"Stars": "1,573",
			"Repo": "https://github.com/vuejs/vite?from=pagser"
		}
	]
}
-------------
```

### Colly Example

Work with colly:
```golang

p := pagser.New()


// On every a element which has href attribute call callback
collector.OnHTML("body", func(e *colly.HTMLElement) {
	//data parser model
	var data PageData
	//parse html data
	err := p.ParseSelection(&data, e.Dom)

})

```

- [See Examples](https://github.com/foolin/pagser/tree/master/_examples)
- [See Tests](https://github.com/foolin/pagser/blob/master/pagser_test.go)

## Dependencies
- github.com/PuerkitoBio/goquery
- github.com/spf13/cast

**Extensions:**
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