package router

import (
	customerror "go-dog/error"
	"go-dog/header"
	"go-dog/plugins"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
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
	codec   plugins.Codec
	methods map[string]*methodstruct
}

//NewRouter 创建路由
func NewRouter(codec plugins.Codec) *Router {
	return &Router{
		methods: make(map[string]*methodstruct),
		codec:   codec,
	}
}

//RegisterByMethod 注册方法
func (pointer *Router) RegisterByMethod(name string, fn interface{}) (arg interface{}, reply interface{}) {
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
		panic("第一个参数必须为go-dog/context")
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
	return pointer.new(argType), pointer.new(replyType)
}

//GetMethodArg 获取方法请求的参数
func (pointer *Router) GetMethodArg(method string) interface{} {
	if vali, ok := pointer.methods[method]; ok {
		return pointer.new(vali.argType)
	}
	return nil
}

//Call 调用方法
func (pointer *Router) Call(ctx plugins.Context, req *header.Request) ([]byte, error) {
	val, ok := pointer.methods[strings.ToLower(req.Method)]
	if ok {
		argv := pointer.new(val.argType)
		err := pointer.codec.DeCode(req.Arg, argv)
		if err != nil {
			return nil, customerror.EnCodeError(customerror.ParamError, "参数不合法")
		}
		if val.argType.Kind() != reflect.Ptr {
			returnValues := val.method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(argv).Elem()})
			errInter := returnValues[1].Interface()
			if errInter != nil {
				return nil, errInter.(error)
			}
			back := returnValues[0].Interface()
			reply, err := pointer.codec.EnCode(back)
			if err != nil {
				return nil, customerror.EnCodeError(customerror.ParamError, "返回参数不合法")
			}
			return reply, nil
		}
		returnValues := val.method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(argv)})
		errInter := returnValues[1].Interface()
		if errInter != nil {
			return nil, errInter.(error)
		}
		back := returnValues[0].Interface()
		reply, err := pointer.codec.EnCode(back)
		if err != nil {
			return nil, customerror.EnCodeError(customerror.InternalServerError, "返回参数不合法")
		}
		return reply, nil

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
