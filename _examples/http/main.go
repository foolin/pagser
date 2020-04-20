package main

import (
	"encoding/json"
	"github.com/foolin/pagser"
	"log"
	"net/http"
)

type ExampPage struct {
	Title string `pagser:"title"`
	H1    string `pagser:"h1"`
	Navs  []struct {
		ID   int    `pagser:"->attrInt(id, -1)"`
		Name string `pagser:"a"`
		Url  string `pagser:"a->attr(href)"`
	} `pagser:".navlink li"`
}

func main() {
	resp, err := http.Get("https://raw.githubusercontent.com/foolin/pagser/master/_examples/pages/demo.html")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//New default config
	p := pagser.New()
	//data parser model
	var page ExampPage
	//parse html data
	err = p.ParseReader(&page, resp.Body)
	//check error
	if err != nil {
		panic(err)
	}

	//print data
	log.Printf("Page data json: \n-------------\n%v\n-------------\n", toJson(page))
}

func toJson(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}
