package pagser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cast"
	"io"
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

// Parse parse html to struct
func (p *Pagser) ParseReader(v interface{}, reader io.Reader) (err error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return err
	}
	return p.ParseDocument(v, doc)
}

// ParseDocument parse document to struct
func (p *Pagser) ParseDocument(v interface{}, document *goquery.Document) (err error) {
	return p.ParseSelection(v, document.Selection)
}

func (p *Pagser) ParseSelection(v interface{}, selection *goquery.Selection) (err error) {
	return p.doParse(v, nil, selection)
}

// ParseSelection parse selection to struct
func (p *Pagser) doParse(v interface{}, stackRefValues []reflect.Value, selection *goquery.Selection) (err error) {
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

		var callOutValue interface{}
		var callErr error
		if tager.FuncName != "" {
			callOutValue, callErr = p.findAndExecFunc(objRefValue, stackRefValues, tager, node)
			if callErr != nil {
				return fmt.Errorf("tag=`%v` parse func error: %v", tagValue, callErr)
			}
			svErr := setRefectValue(fieldType.Type.Kind(), fieldValue, callOutValue)
			if svErr != nil {
				return fmt.Errorf("tag=`%v` set value error: %v", tagValue, svErr)
			}

			//goto parse next field
			continue
		}

		if stackRefValues == nil {
			stackRefValues = make([]reflect.Value, 0)
		}
		stackRefValues = append(stackRefValues, objRefValue)

		//set value
		kind := fieldType.Type.Kind()
		switch {
		case kind == reflect.Ptr:
			subModel := reflect.New(fieldType.Type.Elem())
			fieldValue.Set(subModel)
			err = p.doParse(subModel.Interface(), stackRefValues, node)
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
					err = p.doParse(itemValue.Addr().Interface(), stackRefValues, subNode)
					if err != nil {
						err = fmt.Errorf("tag=`%v` %#v parser error: %v", tagValue, itemValue, err)
						return false
					}
				case itemKind == reflect.Ptr && itemValue.Type().Elem().Kind() == reflect.Struct:
					itemValue = reflect.New(itemType.Elem())
					err = p.doParse(itemValue.Interface(), stackRefValues, subNode)
					if err != nil {
						err = fmt.Errorf("tag=`%v` %#v parser error: %v", tagValue, itemValue, err)
						return false
					}
				default:
					//slice.Index(i).Set(itemValue)
					if tager.FuncName != "" {
						itemOutValue, itemErr := p.findAndExecFunc(objRefValue, stackRefValues, tager, subNode)
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
			err = p.doParse(subModel.Interface(), stackRefValues, node)
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
			fieldValue.SetString(strings.TrimSpace(node.Text()))
		}
	}
	return nil
}

/**
fieldType := refTypeElem.Field(i)
fieldValue := refValueElem.Field(i)
*/
func (p *Pagser) findAndExecFunc(objRefValue reflect.Value, stackRefValues []reflect.Value, selTag *Tager, node *goquery.Selection) (interface{}, error) {
	if selTag.FuncName != "" {

		//call object method
		callMethod := findMethod(objRefValue, selTag.FuncName)
		if callMethod.IsValid() {
			//execute method
			return execMethod(callMethod, selTag, node)
		}

		//call root method
		size := len(stackRefValues)
		if size > 0 {
			for i := size - 1; i >= 0; i-- {
				callMethod = findMethod(stackRefValues[i], selTag.FuncName)
				if callMethod.IsValid() {
					//execute method
					return execMethod(callMethod, selTag, node)
				}
			}
		}

		//global function
		if fn, ok := p.funcs[selTag.FuncName]; ok {
			outValue, err := fn(node, selTag.FuncParams...)
			if err != nil {
				return nil, fmt.Errorf("call registered func %v error: %v", selTag.FuncName, err)
			}
			return outValue, nil
		}

		//not found method
		return nil, fmt.Errorf("not found method %v", selTag.FuncName)
	}
	return strings.TrimSpace(node.Text()), nil
}

func findMethod(objRefValue reflect.Value, funcName string) reflect.Value {
	callMethod := objRefValue.MethodByName(funcName)
	if callMethod.IsValid() {
		return callMethod
	}
	//call element method
	return objRefValue.Elem().MethodByName(funcName)
}

func execMethod(callMethod reflect.Value, selTag *Tager, node *goquery.Selection) (interface{}, error) {
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
	case kind == reflect.String:
		kv, err := cast.ToStringE(v)
		if err != nil {
			return err
		}
		fieldValue.SetString(kv)
	//case kind == reflect.Interface:
	//	fieldValue.Set(reflect.ValueOf(v))
	default:
		fieldValue.Set(reflect.ValueOf(v))
		//return fmt.Errorf("not support type %v", kind)
	}
	return nil
}
