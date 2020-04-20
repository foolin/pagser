package pagser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cast"
	"reflect"
	"strings"
)

// Parse parse html to struct
func (p *Pagser) Parse(v interface{}, document string) (err error) {
	reader, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		return err
	}
	return p.ParseDocument(v, reader)
}

// ParseDocument parse document to struct
func (p *Pagser) ParseDocument(v interface{}, document *goquery.Document) (err error) {
	return p.ParseSelection(v, document.Selection)
}

// ParseSelection parse selection to struct
func (p *Pagser) ParseSelection(v interface{}, selection *goquery.Selection) (err error) {
	objRefType := reflect.TypeOf(v)
	objRefValue := reflect.ValueOf(v)
	//log.Printf("%#v kind is %v | %v", v, objRefValue.Kind(), reflect.Ptr)
	if objRefValue.Kind() != reflect.Ptr {
		return fmt.Errorf("%v is non-pointer", objRefType)
	}
	if objRefValue.IsNil() {
		return fmt.Errorf("%v is nil", objRefType)
	}
	objRefTypeElem := objRefType.Elem()
	objRefValueElem := objRefValue.Elem()

	for i := 0; i < objRefValueElem.NumField(); i++ {
		fieldType := objRefTypeElem.Field(i)
		fieldValue := objRefValueElem.Field(i)
		//tagValue := fieldType.Tag.Get(parserTagName)
		tagValue, tagOk := fieldType.Tag.Lookup(p.config.TagerName)
		if !tagOk {
			if p.config.Debug {
				fmt.Printf("[WARN] not found tager name=[%v] in field: %v, eg: `%v:\".navlink a->attr(href)\"`\n",
					p.config.TagerName, fieldType.Name, p.config.TagerName)
			}
			continue
		}
		if tagValue == p.config.IgnoreSymbol {
			continue
		}

		tager, ok := p.tagers[tagValue]
		if !ok || tager == nil {
			tager = p.newTager(tagValue)
			p.tagers[tagValue] = tager
		}

		node := selection
		if tager.Selector != "" {
			node = selection.Find(tager.Selector)
		}

		//set value
		kind := fieldType.Type.Kind()
		switch {
		case kind == reflect.Ptr:
			subModel := reflect.New(fieldType.Type.Elem())
			fieldValue.Set(subModel)
			err = p.ParseSelection(subModel.Interface(), node)
			if err != nil {
				return fmt.Errorf("tag=`%v` %#v parser error: %v", tagValue, subModel, err)
			}
			//Slice
		case kind == reflect.Slice:
			sliceType := fieldValue.Type()
			itemType := sliceType.Elem()
			itemKind := itemType.Kind()
			slice := reflect.MakeSlice(sliceType, node.Size(), node.Size())
			node.EachWithBreak(func(i int, subNode *goquery.Selection) bool {
				//outhtml, _ := goquery.OuterHtml(subNode)
				//log.Printf("%v => %v", i, outhtml)
				itemValue := reflect.New(itemType).Elem()
				switch {
				case itemKind == reflect.Struct:
					err = p.ParseSelection(itemValue.Addr().Interface(), subNode)
					if err != nil {
						err = fmt.Errorf("tag=`%v` %#v parser error: %v", tagValue, itemValue, err)
						return false
					}
				case itemKind == reflect.Ptr && itemValue.Type().Elem().Kind() == reflect.Struct:
					itemValue = reflect.New(itemType.Elem())
					err = p.ParseSelection(itemValue.Interface(), subNode)
					if err != nil {
						err = fmt.Errorf("tag=`%v` %#v parser error: %v", tagValue, itemValue, err)
						return false
					}
				default:
					//slice.Index(i).Set(itemValue)
					if tager.FuncName != "" {
						itemOutValue, itemErr := p.callFuncFieldValue(objRefValue, tager, subNode)
						//fmt.Printf("call slice func %v value: %v\n", tager.FuncName, itemOutValue)
						if itemErr != nil {
							err = fmt.Errorf("tag=`%v` parse slice item error: %v", tagValue, itemErr)
							return false
						}
						svErr := setRefectValue(itemType.Kind(), itemValue, itemOutValue)
						if err != nil {
							err = fmt.Errorf("tag=`%v` set value error: %v", tagValue, svErr)
							return false
						}
					} else {
						itemValue.SetString(strings.TrimSpace(subNode.Text()))
					}
				}
				slice.Index(i).Set(itemValue)
				return true
			})
			if err != nil {
				return err
			}
			fieldValue.Set(slice)
		case kind == reflect.Struct:
			subModel := reflect.New(fieldType.Type)
			err = p.ParseSelection(subModel.Interface(), node)
			if err != nil {
				return fmt.Errorf("tag=`%v` %#v parser error: %v", tagValue, subModel, err)
			}
			fieldValue.Set(subModel.Elem())
			//UnsafePointer
			//Complex64
			//Complex128
			//Array
			//Chan
			//Func
		default:
			if tager.FuncName != "" {
				callOutValue, callErr := p.callFuncFieldValue(objRefValue, tager, node)
				if callErr != nil {
					return fmt.Errorf("tag=`%v` parse func error: %v", tagValue, callErr)
				}
				//fmt.Printf("call func %v value: %#v\n", tager.FuncName, callOutValue)
				//fieldValue.Set(reflect.ValueOf(callOutValue))
				svErr := setRefectValue(fieldType.Type.Kind(), fieldValue, callOutValue)
				if svErr != nil {
					return fmt.Errorf("tag=`%v` set value error: %v", tagValue, svErr)
				}
			} else {
				fieldValue.SetString(strings.TrimSpace(node.Text()))
			}
		}
	}
	return nil
}

/**
fieldType := refTypeElem.Field(i)
fieldValue := refValueElem.Field(i)
*/
func (p *Pagser) callFuncFieldValue(objRefValue reflect.Value, selTag *Tager, node *goquery.Selection) (interface{}, error) {
	if selTag.FuncName != "" {
		if fn, ok := p.funcs[selTag.FuncName]; ok {
			outValue, err := fn(node, selTag.FuncParams...)
			if err != nil {
				return nil, fmt.Errorf("call registered func %v error: %v", selTag.FuncName, err)
			}
			return outValue, nil
		}

		refValueElem := objRefValue.Elem()
		callMethod := refValueElem.MethodByName(selTag.FuncName)
		if !callMethod.IsValid() {
			callMethod = objRefValue.MethodByName(selTag.FuncName)
		}
		if !callMethod.IsValid() {
			return nil, fmt.Errorf("not found method %v", selTag.FuncName)
		}
		callParams := make([]reflect.Value, 0)
		callParams = append(callParams, reflect.ValueOf(node))

		callReturns := callMethod.Call(callParams)
		if len(callReturns) <= 0 {
			return nil, fmt.Errorf("method %v not return any value", selTag.FuncName)
		}
		if len(callReturns) > 1 {
			if err, ok := callReturns[len(callReturns)-1].Interface().(error); ok {
				if err != nil {
					return nil, fmt.Errorf("method %v return error: %v", selTag.FuncName, err)
				}
			}
		}
		return callReturns[0].Interface(), nil
	}
	return strings.TrimSpace(node.Text()), nil
}

func setRefectValue(kind reflect.Kind, fieldValue reflect.Value, v interface{}) (err error) {
	//set value
	switch {
	//Bool
	case kind == reflect.Bool:
		kv, err := cast.ToBoolE(v)
		if err != nil {
			return err
		}
		fieldValue.SetBool(kv)
	case kind >= reflect.Int && kind <= reflect.Int64:
		kv, err := cast.ToInt64E(v)
		if err != nil {
			return err
		}
		fieldValue.SetInt(kv)
	case kind >= reflect.Uint && kind <= reflect.Uintptr:
		kv, err := cast.ToUint64E(v)
		if err != nil {
			return err
		}
		fieldValue.SetUint(kv)
		//Float32
		//Float64
	case kind == reflect.Float32 || kind == reflect.Float64:
		value, err := cast.ToFloat64E(v)
		if err != nil {
			return err
		}
		fieldValue.SetFloat(value)
		//Interface
	case kind == reflect.Interface:
		fieldValue.Set(reflect.ValueOf(v))
	case kind == reflect.String:
		kv, err := cast.ToStringE(v)
		if err != nil {
			return err
		}
		fieldValue.SetString(kv)
	default:
		return fmt.Errorf("not support type %v", kind)
	}
	return nil
}
