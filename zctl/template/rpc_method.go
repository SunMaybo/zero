package template

const RPCMethodTemplate = `package {{.PackageName}}

import "github.com/SunMaybo/zero/common/zrpc"

const (
{{range $index, $method := .Names}}
	Rpc_{{$method}} zrpc.RpcMethod = "{{$method}}"
{{end}}
)
`

type MethodName struct {
	Names       []string
	PackageName string
}
