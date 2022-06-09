package template

const JavaRPCImplPattern = `package {{.PackageName}};

import io.grpc.stub.StreamObserver;


public abstract class {{.ServiceName}} extends {{.GrpcFileName}}.{{.ServiceBaseName}} {
   {{range $index, $method := .MethodSigns}}
	{{if eq $method.IsStream 0}}
    @Override
    public void {{$method.ReturnParam}} {{$method.Method}}({{$method.Param1}}, {{$method.Param2}}) {
 		try {        
			responseObserver.onNext({{$method.Method}}(request));
        	responseObserver.onCompleted();
        } catch (Exception e) {
            responseObserver.onError(e);
        }
    }
    {{$method.MethodComment}}
    protected abstract {{$method.Param2T}} {{$method.Method}}({{$method.Param1}});
   {{else if eq $method.IsStream 1}}
    @Override
    public {{$method.ReturnParam}} {{$method.Method}}({{$method.Param1}}) {
        return {{$method.Method}}WithDuplex(responseObserver);
    }
    {{$method.MethodComment}}
    protected abstract {{$method.ReturnParam}} {{$method.Method}}WithDuplex({{$method.Param1}});

   {{else}}
    @Override
    public void {{$method.Method}}({{$method.Param1}}, {{$method.Param2}}) {
        {{$method.Method}}WithReturnSimplex(request, responseObserver);
    }
	{{$method.MethodComment}}
    public abstract void {{$method.Method}}WithReturnSimplex({{$method.Param1}}, {{$method.Param2}});
   {{end}}
   {{end}}
}
`
const (
	IsNotStream MethodType = iota
	DuplexStream
	SimplexStream
	ReturnSimplexStream
)

type MethodType int

type JavaRpcImpl struct {
	PackageName     string
	ServiceName     string
	GrpcFileName    string
	ServiceBaseName string
	MethodSigns     []Method
}

type Method struct {
	Method        string
	Param1        string
	Param2        string
	Param2T       string
	ReturnParam   string
	MethodComment string
	IsStream      MethodType
}
