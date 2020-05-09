package pagser

import (
	"testing"
)

//test errors
func TestBuiltinSelectionsErrors(t *testing.T) {
	tests := []funcWantError{
		//not args
		{true, "eq", []string{}, ``},
		//index not number
		{true, "eq", []string{"a"}, ``},
		//not args
		{true, "parentsUntil", []string{}, ``},
	}

	for _, tt := range tests {
		var sel = newTewSelection(tt.data)
		_, err := builtinFuncs[tt.fun](sel, tt.args...)
		if tt.want {
			if err == nil {
				t.Errorf("%v want an error", tt.String())
			}
			continue
		}
		if err != nil {
			t.Errorf("%v want no error, but error is %v", tt.String(), err)
		}
	}
}
