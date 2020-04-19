package ugchtml

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/foolin/pagser"
	"github.com/microcosm-cc/bluemonday"
)

func UgcHtml(node *goquery.Selection, args ...string) (interface{}, error) {
	html, err := goquery.OuterHtml(node)
	if err != nil {
		return html, err
	}
	p := bluemonday.UGCPolicy()
	// The policy can then be used to sanitize lots of input and it is safe to use the policy in multiple goroutines
	return p.Sanitize(html), nil
}

func Register(p *pagser.Pagser) {
	p.RegisterFunc("UgcHtml", UgcHtml)
}
