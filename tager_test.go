package pagser

import (
	"fmt"
	"testing"
)

func TestParseFuncParams(t *testing.T) {
	//const str = `"Hello world, Andy", abc, 123, 1.2,3, -1, ',','abc\'x', -12.3, true, '1,2,3', false, "[1,2,3]",bool`
	//const str = `Foo bar random "letters lol" stuff`
	const str = `href,123, 456, 'Foolin said: "it\'s ok"', "not work", 'Hello world'`
	//const str = "href"
	params := parseFuncParams(str)
	for i, v := range params {
		fmt.Printf("%.2d: %v\n", i, v)
	}
	fmt.Printf("%v", prettyJson(params))
}
