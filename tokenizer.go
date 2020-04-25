package pagser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type tokenState int

const (
	tokenInput tokenState = iota
	tokenStartQuote
	tokenEndQoute
	tokenComma
)

//->fn()
//->fn(xxx)
//->fn('xxx')
//->fn('xxx\'xxx', 'xxx,xxx')
var rxFunc = regexp.MustCompile("^\\s*([a-zA-Z]+)\\s*(\\(([^\\)]*)\\))?\\s*$")

// tagTokenizer struct tag info
type tagTokenizer struct {
	Selector   string
	FuncName   string
	FuncParams []string
}

func (p *Pagser) newTag(tagValue string) (*tagTokenizer, error) {
	//fmt.Println("tag value: ", tagValue)
	tag := &tagTokenizer{}
	if tagValue == "" {
		return tag, nil
	}
	selectors := strings.Split(tagValue, p.Config.FuncSymbol)
	funcValue := ""
	for i := 0; i < len(selectors); i++ {
		switch i {
		case 0:
			tag.Selector = strings.TrimSpace(selectors[i])
		case 1:
			funcValue = selectors[i]
		}
	}
	matches := rxFunc.FindStringSubmatch(funcValue)
	if len(matches) < 3 {
		return tag, nil
	}
	tag.FuncName = strings.TrimSpace(matches[1])
	//tag.FuncParams = strings.Split(matches[2], ",")
	params, err := parseFuncParamTokens(matches[3])
	if err != nil {
		return nil, fmt.Errorf("tag=`%v` is invalid: %v", tagValue, err)
	}
	tag.FuncParams = params
	if p.Config.Debug {
		fmt.Printf("----- debug -----\n`%v`\n%v\n", tagValue, prettyJson(tag))
	}
	return tag, nil
}

func parseFuncParamTokens(text string) ([]string, error) {
	tokens := make([]string, 0)
	textLen := len(text)
	token := strings.Builder{}

	var currentState tokenState
	for pos := 0; pos < textLen; pos++ {
		ch := rune(text[pos])
		switch ch {
		case '\'':
			if currentState == tokenStartQuote {
				tokens = append(tokens, token.String())
				token.Reset()
				currentState = tokenEndQoute
				continue
			}
			if token.Len() <= 0 {
				currentState = tokenStartQuote
				continue
			}
		case ',':
			if currentState == tokenStartQuote {
				token.WriteRune(ch)
				continue
			}
			if currentState == tokenComma || token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
				currentState = tokenComma
			}
			continue
		case '\\':
			if currentState == tokenStartQuote && pos+1 < textLen {
				//escape character -> "\'"
				nextCh := rune(text[pos+1])
				if nextCh == '\'' {
					token.WriteRune(nextCh)
					pos += 1
					continue
				}
			}
		case ' ':
			if (currentState != tokenStartQuote) && token.Len() <= 0 {
				continue
			}
		}

		token.WriteRune(ch)
	}
	if currentState == tokenStartQuote {
		return []string{}, errors.New("syntax error, single quote not closed")
	}
	//end
	if token.Len() > 0 {
		tokens = append(tokens, token.String())
		token.Reset()
	}
	return tokens, nil
}
