package pagser

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cast"
)

// Parse parse html to struct
func (p *Pagser) Parse(v interface{}, document string) (err error) {
	reader, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		return err
	}
	return p.ParseDocument(v, reader)
}

// ParseReader parse html to struct
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

// ParseSelection parse selection to struct
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
		kind := fieldType.Type.Kind()

		//tagValue := fieldType.Tag.Get(parserTagName)
		tagValue, tagOk := fieldType.Tag.Lookup(p.Config.TagName)
		if !tagOk {
			if p.Config.Debug {
				fmt.Printf("[INFO] not found tag name=[%v] in field: %v, eg: `%v:\".navlink a->attr(href)\"`\n",
					p.Config.TagName, fieldType.Name, p.Config.TagName)
			}
			continue
		}
		if tagValue == ignoreSymbol {
			continue
		}

		cacheTag, ok := p.mapTags.Load(tagValue)
		var tag *tagTokenizer
		if !ok || cacheTag == nil {
			tag, err = p.newTag(tagValue)
			if err != nil {
				return err
			}
			p.mapTags.Store(tagValue, tag)
		} else {
			tag = cacheTag.(*tagTokenizer)
		}

		node := selection
		if tag.Selector != "" {
			node = selection.Find(tag.Selector)
		}

		var callOutValue interface{}
		var callErr error
		if tag.FuncName != "" {
			callOutValue, callErr = p.findAndExecFunc(objRefValue, stackRefValues, tag, node)
			if callErr != nil {
				return fmt.Errorf("tag=`%v` parse func error: %v", tagValue, callErr)
			}
			if subNode, ok := callOutValue.(*goquery.Selection); ok {
				//set sub node to current node
				node = subNode
			} else {
				svErr := p.setRefectValue(fieldType.Type.Kind(), fieldValue, callOutValue)
				if svErr != nil {
					return fmt.Errorf("tag=`%v` set value error: %v", tagValue, svErr)
				}
				//goto parse next field
				continue
			}
		}

		if stackRefValues == nil {
			stackRefValues = make([]reflect.Value, 0)
		}
		stackRefValues = append(stackRefValues, objRefValue)

		//set value
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
					itemValue.SetString(strings.TrimSpace(subNode.Text()))
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
func (p *Pagser) findAndExecFunc(objRefValue reflect.Value, stackRefValues []reflect.Value, selTag *tagTokenizer, node *goquery.Selection) (interface{}, error) {
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
		if fn, ok := p.mapFuncs.Load(selTag.FuncName); ok {
			cfn := fn.(CallFunc)
			outValue, err := cfn(node, selTag.FuncParams...)
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

func execMethod(callMethod reflect.Value, selTag *tagTokenizer, node *goquery.Selection) (interface{}, error) {
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

func (p Pagser) setRefectValue(kind reflect.Kind, fieldValue reflect.Value, v interface{}) (err error) {
	//set value
	switch {
	//Bool
	case kind == reflect.Bool:
		if p.Config.CastError {
			kv, err := cast.ToBoolE(v)
			if err != nil {
				return err
			}
			fieldValue.SetBool(kv)
		} else {
			fieldValue.SetBool(cast.ToBool(v))
		}
	case kind >= reflect.Int && kind <= reflect.Int64:
		if p.Config.CastError {
			kv, err := cast.ToInt64E(v)
			if err != nil {
				return err
			}
			fieldValue.SetInt(kv)
		} else {
			fieldValue.SetInt(cast.ToInt64(v))
		}
	case kind >= reflect.Uint && kind <= reflect.Uintptr:
		if p.Config.CastError {
			kv, err := cast.ToUint64E(v)
			if err != nil {
				return err
			}
			fieldValue.SetUint(kv)
		} else {
			fieldValue.SetUint(cast.ToUint64(v))
		}
	case kind == reflect.Float32 || kind == reflect.Float64:
		if p.Config.CastError {
			value, err := cast.ToFloat64E(v)
			if err != nil {
				return err
			}
			fieldValue.SetFloat(value)
		} else {
			fieldValue.SetFloat(cast.ToFloat64(v))
		}
	case kind == reflect.String:
		if p.Config.CastError {
			kv, err := cast.ToStringE(v)
			if err != nil {
				return err
			}
			fieldValue.SetString(kv)
		} else {
			fieldValue.SetString(cast.ToString(v))
		}
	case kind == reflect.Slice || kind == reflect.Array:
		sliceType := fieldValue.Type().Elem()
		itemKind := sliceType.Kind()
		if p.Config.CastError {
			switch itemKind {
			case reflect.Bool:
				kv, err := cast.ToBoolSliceE(v)
				if err != nil {
					return err
				}
				fieldValue.Set(reflect.ValueOf(kv))
			case reflect.Int:
				kv, err := cast.ToIntSliceE(v)
				if err != nil {
					return err
				}
				fieldValue.Set(reflect.ValueOf(kv))
			case reflect.Int32:
				kv, err := toInt32SliceE(v)
				if err != nil {
					return err
				}
				fieldValue.Set(reflect.ValueOf(kv))
			case reflect.Int64:
				kv, err := toInt64SliceE(v)
				if err != nil {
					return err
				}
				fieldValue.Set(reflect.ValueOf(kv))
			case reflect.Float32:
				kv, err := toFloat32SliceE(v)
				if err != nil {
					return err
				}
				fieldValue.Set(reflect.ValueOf(kv))
			case reflect.Float64:
				kv, err := toFloat64SliceE(v)
				if err != nil {
					return err
				}
				fieldValue.Set(reflect.ValueOf(kv))
			case reflect.String:
				kv, err := cast.ToStringSliceE(v)
				if err != nil {
					return err
				}
				fieldValue.Set(reflect.ValueOf(kv))
			default:
				fieldValue.Set(reflect.ValueOf(v))
			}
		} else {
			switch itemKind {
			case reflect.Bool:
				kv := cast.ToBoolSlice(v)
				fieldValue.Set(reflect.ValueOf(kv))
			case reflect.Int:
				kv := cast.ToIntSlice(v)
				fieldValue.Set(reflect.ValueOf(kv))
			case reflect.Int32:
				kv := toInt32Slice(v)
				fieldValue.Set(reflect.ValueOf(kv))
			case reflect.Int64:
				kv := toInt64Slice(v)
				fieldValue.Set(reflect.ValueOf(kv))
			case reflect.Float32:
				kv := toFloat32Slice(v)
				fieldValue.Set(reflect.ValueOf(kv))
			case reflect.Float64:
				kv := toFloat64Slice(v)
				fieldValue.Set(reflect.ValueOf(kv))
			case reflect.String:
				kv := cast.ToStringSlice(v)
				fieldValue.Set(reflect.ValueOf(kv))
			default:
				fieldValue.Set(reflect.ValueOf(v))
			}
		}
	//case kind == reflect.Interface:
	//	fieldValue.Set(reflect.ValueOf(v))
	default:
		fieldValue.Set(reflect.ValueOf(v))
		//return fmt.Errorf("not support type %v", kind)
	}
	return nil
}
