package template

const RPCMethodTemplate = `package {{.PackageName}}

import "github.com/SunMaybo/zero/common/zrpc/interceptor"

const (
{{range $index, $method := .Names}}
	Rpc_{{$method}} interceptor.RpcMethod = "{{$method}}"
{{end}}
)
`

type MethodName struct {
	Names       []string
	PackageName string
}
