package service

const (
	ReqCount      = "request_count"
	ReqDuration   = "request_duration_seconds"
	ReqSizeBytes  = "request_size_bytes"
	RespSizeBytes = "response_size_bytes"
)

const (
	Method = "method"
	Path   = "path"
	Name   = "name"
	Code   = "code"
)

var labels = []string{Name, Method, Path, Code}
