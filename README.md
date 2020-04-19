# Pagser

Pagser inspired by  <u>pag</u>e par<u>ser</u>, 
an HTML deserialization to struct based on [goquery](https://github.com/PuerkitoBio/goquery) and struct tags。
It's simple, easy, extensible, configurable parsing library from [scrago](https://github.com/foolin/scrago).

## Features

* **Simple** - use golang struct tag syntax.
* **Easy** - easy use for your spider/crawler application.
* **Extensible** - Support for extension functions.
* **Struct tag** - Grammar is simple.
* **Configurable** - Reflect
* **GoQuery/Colly** - Support all [goquery](https://github.com/PuerkitoBio/goquery) framework, such as [go-colly](https://github.com/gocolly/colly).

# Usage



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



