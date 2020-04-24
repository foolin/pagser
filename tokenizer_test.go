package pagser

import "testing"

func TestFuncParamTokens(t *testing.T) {
	inputs := []string{
		`''`,
		`','`,
		`',', ','	`,
		`'\''`,
		`'\'' , '\''`,
		`'', ''`,
		`'',,''`,
		`',', '\''`,
		`href`,
		`'href'`,
		`'href', abc`,
		`'href', 'abc'`,
		`'a,b', abc`,
		`'a,b\'', abc`,
		`'it\'s', abc`,
		`'it\'s', 'abc'`,
		`href, ,' abc', 123.123`,
		`'isShow = [', $value, ]`,
		`' this is words: ', [, $value, ]`,
		`it's ok?, hello , world!'`,
		`'hello, i\'m ok!'`,
		`The name is 'ABC'`,
		`'hello, i\'m error!`,
	}
	for _, v := range inputs {
		args, err := parseFuncParamTokens(v)
		if err != nil {
			t.Logf("X [FAIL] (%v) => %v", v, err)
		} else {
			t.Logf("√ [PASS] (%v) => %v", v, prettyJson(args))
		}

	}
}
