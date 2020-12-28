package jaeger

import (
	"errors"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

//Jaeger 对象
type Jaeger struct {
	closer io.Closer
	cfg    plugins.Cfg
}

//NewJaeger 初始化
func NewJaeger(name string, config plugins.Cfg) *Jaeger {
	if config.GetRunmode() == "trace" && config.GetJaeger() != "" {
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: config.GetJaeger(),
			},
		}
		closer, err := cfg.InitGlobalTracer(name)
		if err != nil {
			log.Panicf("Could not initialize jaeger tracer: %s", err.Error())
			return nil
		}
		return &Jaeger{closer: closer, cfg: config}
	}
	return &Jaeger{cfg: config}
}

// StartSpan 开启
func (j *Jaeger) StartSpan(ctx plugins.Context, operationName string) (opentracing.Span, error) {
	if j.cfg.GetRunmode() != "trace" {
		return nil, errors.New("没有开启trace")
	}
	var data string
	var span opentracing.Span
	tracer := opentracing.GlobalTracer()
	if err := ctx.GetDataByKey("opentracing", &data); err == nil {
		carrier := opentracing.TextMapCarrier(map[string]string{"uber-trace-id": data})
		spanContext, err := tracer.Extract(opentracing.TextMap, carrier)
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			log.Errorln(err.Error())
			return nil, err
		}
		span = opentracing.StartSpan(operationName, opentracing.ChildOf(spanContext))
	} else {
		span = opentracing.StartSpan(operationName)
	}
	metadata := make(map[string]string)
	err := tracer.Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(metadata))
	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}
	if err := ctx.SetData("opentracing", metadata["uber-trace-id"]); err != nil {
		log.Errorln(err.Error())
	}
	return span, nil
}

//Request 请求
func (j *Jaeger) Request(ctx plugins.Context, servicename, method string, request interface{}) {
	if j.cfg.GetRunmode() == "trace" {
		span, err := j.StartSpan(ctx, servicename+"."+method)
		if err == nil {
			ctx.SetShare("Span", span)
		}
	}

}

//Respone 响应
func (j *Jaeger) Respone(ctx plugins.Context, servicename, method string, respone interface{}, err error) {
	if j.cfg.GetRunmode() == "trace" {
		if span, ok := ctx.GetShareByKey("Span").(opentracing.Span); ok {
			span.Finish()
		}
	}
}

//Close 关闭
func (j *Jaeger) Close() error {
	if j.closer != nil {
		return j.closer.Close()
	}
	return nil
}
