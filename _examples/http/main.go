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
		Name        string `pagser:"h1"`
		Description string `pagser:"h1 + p"`
		Stars       string `pagser:"a.muted-link->eq(0)"`
		Repo        string `pagser:"h1 a->concatAttr(href, 'github.com', $value)"`
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
