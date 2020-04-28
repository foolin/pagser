package pagser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

// Builtin functions are registered with a lowercase initial, eg: Text -> text()
type BuiltinSelections struct {
}

// child(selector='') gets the child elements of each element in the Selection,
// Filtered by the specified selector if selector not empty,
// It returns Selection object containing these elements for nested struct..
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->child()"`
//	}
func (builtin BuiltinSelections) Child(node *goquery.Selection, args ...string) (out interface{}, err error) {
	selector := ""
	if len(args) > 0 {
		selector = strings.TrimSpace(args[0])
	}
	if selector != "" {
		return node.ChildrenFiltered(selector), nil
	}
	return node.Children(), nil
}

// eq(index) reduces the set of matched elements to the one at the specified index.
// If a negative index is given, it counts backwards starting at the end of the set.
// It returns a Selection object for nested struct, and an empty Selection object if the
// index is invalid.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->eq(0)"`
//	}
func (builtin BuiltinSelections) Eq(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("nodeEq(index) must has `index` value")
	}
	indexValue := strings.TrimSpace(args[0])
	idx, err := strconv.Atoi(indexValue)
	if err != nil {
		return "", fmt.Errorf("index=`" + indexValue + "` is not number: " + err.Error())
	}
	return node.Eq(idx), nil
}

// first() First reduces the set of matched elements to the first in the set.
// It returns a new Selection object, and an empty Selection object if the
// the selection is empty.
// It returns Selection object containing these elements for nested struct.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->first()"`
//	}
func (builtin BuiltinSelections) First(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return node.First(), nil
}

// last(selector='') reduces the set of matched elements to the last in the set.
// It returns a new Selection object, and an empty Selection object if
// the selection is empty.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->last()"`
//	}
func (builtin BuiltinSelections) Last(node *goquery.Selection, args ...string) (out interface{}, err error) {
	return node.Last(), nil
}

// next(selector='') gets the immediately following sibling of each element in the Selection.
// Filtered by the specified selector if selector not empty,
// It returns Selection object containing these elements for nested struct.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->next()"`
//	}
func (builtin BuiltinSelections) Next(node *goquery.Selection, args ...string) (out interface{}, err error) {
	selector := ""
	if len(args) > 0 {
		selector = strings.TrimSpace(args[0])
	}
	if selector != "" {
		return node.NextFiltered(selector), nil
	}
	return node.Next(), nil
}

// parent(selector='') gets the parent elements of each element in the Selection.
// Filtered by the specified selector if selector not empty,
// It returns Selection object containing these elements for nested struct.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->parent()"`
//	}
func (builtin BuiltinSelections) Parent(node *goquery.Selection, args ...string) (out interface{}, err error) {
	selector := ""
	if len(args) > 0 {
		selector = strings.TrimSpace(args[0])
	}
	if selector != "" {
		return node.ParentFiltered(selector), nil
	}
	return node.Parent(), nil
}

// parents(selector='') gets the parent elements of each element in the Selection.
// Filtered by the specified selector if selector not empty,
// It returns Selection object containing these elements for nested struct.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->parents()"`
//	}
func (builtin BuiltinSelections) Parents(node *goquery.Selection, args ...string) (out interface{}, err error) {
	selector := ""
	if len(args) > 0 {
		selector = strings.TrimSpace(args[0])
	}
	if selector != "" {
		return node.ParentsFiltered(selector), nil
	}
	return node.Parents(), nil
}

// parentsUntil(selector) gets the ancestors of each element in the Selection, up to but
// not including the element matched by the selector. It returns a new Selection
// object containing the matched elements.
// It returns Selection object containing these elements for nested struct.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->parentsUntil('.wrap')"`
//	}
func (builtin BuiltinSelections) ParentsUntil(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("parentsUntil must has selector")
	}
	selector := strings.TrimSpace(args[0])
	return node.ParentsUntil(selector), nil
}

// prev() gets the immediately preceding sibling of each element in the Selection.
// Filtered by the specified selector if selector not empty,
// It returns Selection object containing these elements for nested struct.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->prev()"`
//	}
func (builtin BuiltinSelections) Prev(node *goquery.Selection, args ...string) (out interface{}, err error) {
	selector := ""
	if len(args) > 0 {
		selector = strings.TrimSpace(args[0])
	}
	if selector != "" {
		return node.PrevFiltered(selector), nil
	}
	return node.Prev(), nil
}

// siblings() gets the siblings of each element in the Selection.
// Filtered by the specified selector if selector not empty,
// It returns Selection object containing these elements for nested struct.
//	struct {
//		SubStruct struct {
//			Example string `pagser:".selector->text()"`
//		}	`pagser:".selector->siblings()"`
//	}
func (builtin BuiltinSelections) Siblings(node *goquery.Selection, args ...string) (out interface{}, err error) {
	selector := ""
	if len(args) > 0 {
		selector = strings.TrimSpace(args[0])
	}
	if selector != "" {
		return node.SiblingsFiltered(selector), nil
	}
	return node.Siblings(), nil
}
