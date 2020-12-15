package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/serviceinfo"
)

//Docs 文档内容
type Docs struct {
	Swagger     string                 `json:"swagger"`
	Info        Info                   `json:"info"`
	Host        string                 `json:"host"`
	BasePath    string                 `json:"basePath"`
	Paths       map[string]interface{} `json:"paths"`
	Definitions map[string]Definitions `json:"definitions"`
}

//Info 信息
type Info struct {
	Description string `json:"description"`
	Title       string `json:"title"`
	Contact     struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"contact"`
	License struct {
	} `json:"license"`
	Version string `json:"version"`
}

//Definitions 参数定义
type Definitions struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Properties map[string]Description `json:"properties"`
}

//Description 描述
type Description struct {
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Items       interface{} `json:"items,omitempty"`
	Ref         string      `json:"$ref,omitempty"`
}

//Ref 链接
type Ref struct {
	Ref string `json:"$ref,omitempty"`
}

//POSTAPI POST API结构体
type POSTAPI struct {
	Post Body `json:"post"`
}

//GETAPI GETAPI API结构体
type GETAPI struct {
	Get Body `json:"get"`
}

//Body 请求
type Body struct {
	Consumes   []string     `json:"consumes"`
	Produces   []string     `json:"produces"`
	Tags       []string     `json:"tags"`
	Summary    string       `json:"summary"`
	Parameters []Parameters `json:"parameters"`
	Responses  struct {
		Code200 struct {
			Description string `json:"description"`
			Schema      struct {
				Type string `json:"type"`
				Ref  string `json:"$ref,omitempty"`
			} `json:"schema"`
		} `json:"200"`
	} `json:"responses"`
}

//Parameters api描述
type Parameters struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description"`
	Name        string `json:"name"`
	In          string `json:"in"`
	Required    bool   `json:"required"`
	Schema      struct {
		Type string `json:"type"`
		Ref  string `json:"$ref,omitempty"`
	} `json:"schema"`
}

//t type解析
func t(tp string) string {
	switch tp {
	case "int8":
		return "integer"
	case "int":
		return "integer"
	case "int32":
		return "integer"
	case "int64":
		return "integer"
	case "uint8":
		return "integer"
	case "uint":
		return "integer"
	case "uint32":
		return "integer"
	case "uint64":
		return "integer"
	case "float":
		return "number"
	case "float32":
		return "number"
	case "float64":
		return "number"
	case "byte":
		return "string"
	case "bool":
		return "boolean"
	default:
		return tp
	}
}

//transformation 转换
func transformation(tp string, value string) (interface{}, error) {
	switch tp {
	case "int8":
		i, e := strconv.ParseInt(value, 10, 8)
		if e != nil {
			return nil, fmt.Errorf("需要参数是int8 %s是", e.Error())
		}
		return int8(i), nil
	case "int":
		return strconv.Atoi(value)
	case "int32":
		i, e := strconv.ParseInt(value, 10, 32)
		if e != nil {
			return nil, fmt.Errorf("需要参数是int32 %s是", e.Error())
		}
		return int32(i), nil
	case "int64":
		return strconv.ParseInt(value, 10, 64)
	case "uint8":
		i, e := strconv.ParseInt(value, 10, 8)
		if e != nil {
			return nil, e
		}
		return uint8(i), nil
	case "uint":
		i, e := strconv.Atoi(value)
		if e != nil {
			return nil, fmt.Errorf("需要参数是uint %s是", e.Error())
		}
		return uint(i), nil
	case "uint32":
		i, e := strconv.ParseInt(value, 10, 32)
		if e != nil {
			return nil, fmt.Errorf("需要参数是uint32 %s是", e.Error())
		}
		return uint32(i), nil
	case "uint64":
		i, e := strconv.ParseInt(value, 10, 64)
		if e != nil {
			return nil, fmt.Errorf("需要参数是uint64 %s是", e.Error())
		}
		return uint64(i), nil
	case "float32":
		return strconv.ParseFloat(value, 32)
	case "float64":
		return strconv.ParseFloat(value, 64)
	case "bool":
		return strconv.ParseBool(value)
	case "string":
		return value, nil
	default:
		return tp, fmt.Errorf("暂时不支持此类型参数%s", tp)
	}
}

//createPOSTAPI 创建一个POSTAPI
func createPOSTAPI(tags, summary, name string, isAuth bool, request, respone map[string]interface{}) (a POSTAPI, definitions []Definitions) {
	api := POSTAPI{Post: Body{
		Consumes: []string{"application/json"},
		Produces: []string{"application/json"},
		Tags:     []string{tags},
		Summary:  summary,
	}}
	parameters := Parameters{
		Description: "请求内容",
		Name:        "body",
		In:          "body",
		Required:    true,
	}
	requestName := strings.Replace(tags+"."+name+"Request", "/", ".", -1)
	requestProperties := createDefinitions(requestName, request)
	definitions = append(definitions, requestProperties...)

	parameters.Schema.Type = "object"
	parameters.Schema.Ref = "#/definitions/" + requestName
	api.Post.Parameters = []Parameters{
		{
			Type:        "integer",
			Description: "请求超时时间,单位秒",
			Name:        "TimeOut",
			In:          "header",
			Required:    true,
		},
		{
			Type:        "string",
			Description: "链路请求ID",
			Name:        "TraceID",
			In:          "header",
			Required:    true,
		},
		{
			Type:        "boolean",
			Description: "是否是测试请求",
			Name:        "IsTest",
			In:          "header",
			Required:    true,
		},
	}
	if isAuth {
		api.Post.Parameters = append(api.Post.Parameters, Parameters{
			Type:        "string",
			Description: "验证Token",
			Name:        "Token",
			In:          "header",
			Required:    true,
		})
	}
	api.Post.Parameters = append(api.Post.Parameters, parameters)

	responeName := strings.Replace(tags+"."+name+"Respone", "/", ".", -1)
	responeProperties := createDefinitions(responeName, respone)
	definitions = append(definitions, responeProperties...)

	api.Post.Responses.Code200.Description = "请求成功返回参数"
	api.Post.Responses.Code200.Schema.Type = "object"
	api.Post.Responses.Code200.Schema.Ref = "#/definitions/" + responeName

	return api, definitions
}

//createGETAPI 创建一个GETAPI
func createGETAPI(tags, summary, name string, isAuth bool, request, respone map[string]interface{}) (a GETAPI, definitions []Definitions) {
	api := GETAPI{Get: Body{
		Consumes: []string{"application/json"},
		Tags:     []string{tags},
		Summary:  summary,
	}}
	for key, value := range request {
		if vali, ok := value.(map[string]interface{}); ok {
			des, ok1 := vali["description"]
			tp, ok2 := vali["type"]
			if ok1 == true && ok2 == true {
				api.Get.Parameters = append(api.Get.Parameters, Parameters{
					Type:        t(tp.(string)),
					Description: des.(string),
					Name:        key,
					In:          "query",
					Required:    true,
				})
			}
		}
	}
	api.Get.Parameters = append(api.Get.Parameters,
		Parameters{
			Type:        "integer",
			Description: "请求超时时间,单位秒",
			Name:        "TimeOut",
			In:          "header",
			Required:    true,
		},
		Parameters{
			Type:        "string",
			Description: "链路请求ID",
			Name:        "TraceID",
			In:          "header",
			Required:    true,
		},
		Parameters{
			Type:        "boolean",
			Description: "是否是测试请求",
			Name:        "IsTest",
			In:          "header",
			Required:    true,
		})
	if isAuth {
		api.Get.Parameters = append(api.Get.Parameters, Parameters{
			Type:        "string",
			Description: "验证Token",
			Name:        "Token",
			In:          "header",
			Required:    true,
		})
	}

	responeName := strings.Replace(tags+"."+name+"Respone", "/", ".", -1)
	responeProperties := createDefinitions(responeName, respone)
	definitions = append(definitions, responeProperties...)

	api.Get.Responses.Code200.Description = "请求成功返回参数"
	api.Get.Responses.Code200.Schema.Type = "object"
	api.Get.Responses.Code200.Schema.Ref = "#/definitions/" + responeName

	return api, definitions
}

//createDefinitions 生成Definitions
func createDefinitions(name string, mp map[string]interface{}) (definitions []Definitions) {
	properties := make(map[string]Description)
	for key, value := range mp {
		if vali, ok := value.(map[string]interface{}); ok {
			slice, ok := vali["slice"]
			des, ok1 := vali["description"]
			tp, ok2 := vali["type"]
			if ok {
				mp, o := slice.(map[string]interface{})
				if o == true {
					description := Description{}
					if ok1 {
						description.Description = des.(string)
					}
					if ok2 {
						description.Type = t(tp.(string))
					}
					son := name + "." + key
					definitions = append(definitions, createDefinitions(son, mp)...)
					description.Items = &Ref{
						Ref: "#/definitions/" + son,
					}
					properties[key] = description
				} else {
					description := Description{}
					if ok1 {
						description.Description = des.(string)
					}
					if ok2 {
						description.Type = t(tp.(string))
					}
					description.Items = map[string]string{
						"type": t(vali["slice"].(string)),
					}
					properties[key] = description
				}
				continue
			} else if object, ok3 := vali["object"]; ok3 {
				mp, o := object.(map[string]interface{})
				if o == true {
					description := Description{}
					if ok1 {
						description.Description = des.(string)
					}
					description.Type = "object"
					son := name + "." + key
					definitions = append(definitions, createDefinitions(son, mp)...)
					description.Ref = "#/definitions/" + son

					properties[key] = description
					continue
				}
			}
			description := Description{}
			if ok1 {
				description.Description = des.(string)
			}
			if ok2 {
				description.Type = t(tp.(string))
			}
			properties[key] = description
		}
	}
	definition := Definitions{
		Name:       name,
		Type:       "object",
		Properties: properties,
	}
	definitions = append(definitions, definition)
	return
}

//swagger info
type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{Schemes: []string{}}

//assembleDocs 组装文档
func (g *Gateway) assembleDocs() string {
	info := Info{
		Description: "",
		Title:       "go-dog网管API文档",
		Version:     "{{.Version}}",
	}
	info.Contact.Name = "有bug请联系电话13688460148"
	info.Contact.URL = "tel:13688460148"

	paths := make(map[string]interface{})
	definitions := make(map[string]Definitions)

	g.discovery.RangeAPI(func(url string, api *serviceinfo.ServcieAPI) {
		if api.Method.Kind == "POST" {
			api, d := createPOSTAPI(
				api.Explain+"["+api.Tags+"]",
				api.Method.Explain,
				api.Method.Name,
				api.Method.IsAuth,
				api.Method.Request,
				api.Method.Response)
			paths[url] = api
			for _, definition := range d {
				definitions[definition.Name] = definition
			}
		}
		if api.Method.Kind == "GET" {
			api, d := createGETAPI(
				api.Explain+"["+api.Tags+"]",
				api.Method.Explain,
				api.Method.Name,
				api.Method.IsAuth,
				api.Method.Request,
				api.Method.Response)
			paths[url] = api
			for _, definition := range d {
				definitions[definition.Name] = definition
			}
		}
	})

	docs := &Docs{
		Swagger:     "2.0",
		Host:        "{{.Host}}",
		BasePath:    "{{.BasePath}}",
		Info:        info,
		Paths:       paths,
		Definitions: definitions,
	}
	buff, _ := json.Marshal(docs)
	return string(buff)
}

//ReadDoc 读取文档
func (g *Gateway) ReadDoc() string {
	docs := g.assembleDocs()
	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(docs)
	if err != nil {
		log.Errorln(err.Error())
		return docs
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		log.Errorln(err.Error())
		return docs
	}
	return tpl.String()
}
