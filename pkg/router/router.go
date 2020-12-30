package router

import (
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/plugins"
)

//定义错误类型
var typeOfError = reflect.TypeOf(new(error)).Elem()

//定义context类型
var typeOfContext = reflect.TypeOf(new(plugins.Context)).Elem()

//子服务处理方法
type methodstruct struct {
	name       string
	method     reflect.Value
	classValue reflect.Value
	ctxType    reflect.Type
	argType    reflect.Type
}

//Router api接口对象
type Router struct {
	methods map[string]*methodstruct
}

//NewRouter 创建路由
func NewRouter() *Router {
	return &Router{
		methods: make(map[string]*methodstruct),
	}
}

//RegisterByMethod 注册方法
func (pointer *Router) RegisterByMethod(name string, fn interface{}) (arg map[string]interface{}, reply map[string]interface{}) {
	if _, ok := pointer.methods[name]; ok {
		panic("此函数名称已经存在")
	}
	method, ok := fn.(reflect.Value)
	if !ok {
		method = reflect.ValueOf(fn)
	}
	if method.Kind() != reflect.Func {
		panic("注册的类型必须是一个函数方法")
	}
	mtype := method.Type()
	//入参判断
	if mtype.NumIn() != 2 {
		panic("注册函数的参数数量不正确")
	}
	ctxType := mtype.In(0)
	if !ctxType.Implements(typeOfContext) {
		panic("第一个参数必须为github.com/tang-go/go-dog/context")
	}
	argType := mtype.In(1)
	//返回值判断
	if mtype.NumOut() != 2 {
		panic("返回值不正确")
	}
	replyType := mtype.Out(0)
	//判断最后一个返回值必须为一个错误
	e := mtype.Out(1)
	if !e.Implements(typeOfError) {
		panic("第二个返回值必须为error")
	}
	pointer.methods[strings.ToLower(name)] = &methodstruct{name: name, method: method, ctxType: ctxType, argType: argType}
	return pointer.analysisStruct(nil, "", pointer.new(argType)), pointer.analysisStruct(nil, "", pointer.new(replyType))
}

//analysisStruct 解析参数
func (pointer *Router) analysisStruct(index *int, name string, class interface{}) map[string]interface{} {
	explain := make(map[string]interface{})
	t := reflect.TypeOf(class)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		tgs := map[string]string{
			"type":        t.Kind().String(),
			"description": t.Kind().String(),
		}
		explain[strings.ToLower(t.Name())] = tgs
		return explain
	}
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		//检测每一个字段的深度
		kind := t.Field(i).Type.Kind()
		if kind == reflect.Struct {
			class := pointer.new(t.Field(i).Type)
			if index == nil {
				i := 0
				index = &i
			}
			tg := pointer.analysisStruct(index, "", class)
			tgs := map[string]interface{}{
				"type":        "object",
				"description": t.Field(i).Tag.Get("description"),
				"object":      tg,
			}
			explain[t.Field(i).Tag.Get("json")] = tgs
			continue
		}
		if kind == reflect.Slice {
			class := pointer.new(t.Field(i).Type.Elem())
			classType := reflect.TypeOf(class)
			if classType.Kind() == reflect.Ptr {
				classType = classType.Elem()
			}
			kind := classType.Kind()
			if kind == reflect.Struct {
				if index == nil {
					i := 0
					index = &i
				}
				tg := make(map[string]interface{})
				if classType.Name() == name {
					*index = *index + 1
					if *index < 2 {
						tg = pointer.analysisStruct(index, classType.Name(), class)
					}
				} else {
					tg = pointer.analysisStruct(index, classType.Name(), class)
				}
				tgs := map[string]interface{}{
					"type":        "array",
					"description": t.Field(i).Tag.Get("description"),
					"slice":       tg,
				}
				explain[t.Field(i).Tag.Get("json")] = tgs

			} else {
				tgs := map[string]interface{}{
					"type":        "array",
					"description": t.Field(i).Tag.Get("description"),
					"slice":       kind.String(),
				}
				explain[t.Field(i).Tag.Get("json")] = tgs
			}
			continue
		}
		tgs := map[string]string{
			"type":        t.Field(i).Tag.Get("type"),
			"description": t.Field(i).Tag.Get("description"),
		}
		explain[t.Field(i).Tag.Get("json")] = tgs
	}
	return explain
}

//GetMethodArg 获取方法请求的参数
func (pointer *Router) GetMethodArg(method string) (interface{}, bool) {
	if vali, ok := pointer.methods[strings.ToLower(method)]; ok {
		return pointer.new(vali.argType), true
	}
	return nil, false
}

//Call 调用方法
func (pointer *Router) Call(ctx plugins.Context, method string, argv interface{}) (interface{}, error) {
	val, ok := pointer.methods[strings.ToLower(method)]
	if ok {
		if val.argType.Kind() != reflect.Ptr {
			returnValues := val.method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(argv).Elem()})
			errInter := returnValues[1].Interface()
			if errInter != nil {
				return nil, errInter.(error)
			}
			back := returnValues[0].Interface()
			return back, nil
		}
		returnValues := val.method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(argv)})
		errInter := returnValues[1].Interface()
		if errInter != nil {
			return nil, errInter.(error)
		}
		back := returnValues[0].Interface()
		return back, nil

	}
	return nil, customerror.EnCodeError(customerror.RPCNotFind, "没有找到RPC函数方法")
}

func (pointer *Router) new(t reflect.Type) interface{} {
	var argv reflect.Value
	if t.Kind() == reflect.Ptr {
		argv = reflect.New(t.Elem())
	} else {
		argv = reflect.New(t)
	}
	return argv.Interface()
}

func (pointer *Router) isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

func (pointer *Router) isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return pointer.isExported(t.Name()) || t.PkgPath() == ""
}
