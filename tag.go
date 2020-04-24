package pagser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/scanner"
)

//->text()
//->value()
//->html()
//->attr(xxx, '')
//->attr('', '')
var rxFunc = regexp.MustCompile("^\\s*([a-zA-Z]+)\\s*(\\(([^\\)]*)\\))?\\s*$")

// parseTag struct tag info
type parseTag struct {
	Selector   string
	FuncName   string
	FuncParams []string
}

func (p *Pagser) newTag(tagValue string) *parseTag {
	//fmt.Println("tag value: ", tagValue)
	cssParser := &parseTag{}
	if tagValue == "" {
		return cssParser
	}
	selectors := strings.Split(tagValue, p.Config.FuncSymbol)
	funcValue := ""
	for i := 0; i < len(selectors); i++ {
		switch i {
		case 0:
			cssParser.Selector = strings.TrimSpace(selectors[i])
		case 1:
			funcValue = selectors[i]
		}
	}
	matches := rxFunc.FindStringSubmatch(funcValue)
	if len(matches) < 3 {
		return cssParser
	}
	cssParser.FuncName = strings.TrimSpace(matches[1])
	//cssParser.FuncParams = strings.Split(matches[2], ",")
	cssParser.FuncParams = parseFuncParams(matches[3])
	if p.Config.Debug {
		fmt.Printf("----- debug -----\n`%v`\n%v\n", tagValue, prettyJson(cssParser))
	}
	return cssParser
}

func parseFuncParams(args string) []string {
	//const str = `"Hello world, Andy", abc, 123, 1.2,3,bool`
	//const str = `Foo bar random "letters lol" stuff`
	var s scanner.Scanner
	s.Init(strings.NewReader(args))
	params := make([]string, 0)

	text := strings.Builder{}
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		tmp := s.TokenText()
		if tmp != "," {
			text.WriteString(tmp)
		} else {
			param := text.String()
			if strings.HasPrefix(param, "'") && strings.HasSuffix(param, "'") {
				param = strings.Trim(param, "'")
			}
			if strings.HasPrefix(param, "\"") && strings.HasSuffix(param, "\"") {
				param = strings.Trim(param, "\"")
			}
			params = append(params, param)
			text.Reset()
		}
	}

	if text.Len() > 0 {
		param := text.String()
		if strings.HasPrefix(param, "'") && strings.HasSuffix(param, "'") {
			param = strings.Trim(param, "'")
		}
		params = append(params, param)
	}
	return params
}

func prettyJson(v interface{}) string {
	bytes, _ := json.MarshalIndent(v, "", "\t")
	return string(bytes)
}
