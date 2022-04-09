package template

import (
	"fmt"
	"github.com/SunMaybo/zero/zctl/parser"
	"io/ioutil"
	"testing"
)

var rpcMetadata parser.RpcMetadata

func init() {
	rpcMetadata, _ = parser.Parser("../../proto/hello/test_services.proto")
}

func TestRPCServer(t *testing.T) {
	temps := ServerTemplateParam{
		Project:     "zero",
		PackageName: rpcMetadata.PackageName,
	}
	for _, sign := range rpcMetadata.MethodSigns {
		ms := MethodSign{}
		if !sign.IsStreamReturnParam && !sign.IsStreamParam {
			ms.MethodName = sign.Name
			ms.Sign = ("(ctx context.Context, in *" + rpcMetadata.PackageName + "." + sign.Param + ")") + ("  (*" + rpcMetadata.PackageName + "." + sign.ReturnParam + ", error)")
			ms.Param = "in"
		} else if !sign.IsStreamParam && sign.IsStreamReturnParam {
			//SayStream(Greeter_SayStreamServer) error
			//SayStream1(Greeter_SayStream1Server) error
			//SayStream2(*HelloRequest, Greeter_SayStream2Server) error
			ms.MethodName = sign.Name
			ms.Sign = "(in *" + rpcMetadata.PackageName + "." + sign.Param + ", stream " + rpcMetadata.PackageName + "." + rpcMetadata.ServiceName + "_" + sign.Name + "Server" + ") error"
			ms.Param = "in,stream"
			ms.ISStream = true
		} else {
			ms.MethodName = sign.Name
			ms.Sign = "(stream " + rpcMetadata.PackageName + "." + rpcMetadata.ServiceName + "_" + sign.Name + "Server" + ") error"
			ms.Param = "stream"
			ms.ISStream = true
		}
		temps.MethodSigns = append(temps.MethodSigns, ms)
	}
	result, err := Parser(RPCServerTemplate, temps)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("./../../services/hello/rpc/server/"+"server.go", []byte(result), 0777)
	if err != nil {
		t.Log(err)
	}
	t.Log(result)
}
func TestRPCLogic(t *testing.T) {
	for _, sign := range rpcMetadata.MethodSigns {
		var result string
		var err error
		if !sign.IsStreamReturnParam && !sign.IsStreamParam {
			result, err = Parser(RPCLogicTemplate, LogicTemplateParam{
				PackageName: rpcMetadata.PackageName,
				MethodName:  sign.Name,
				Project:     "zero",
				Sign:        ("(in *" + rpcMetadata.PackageName + "." + sign.Param + ")") + (" (*" + rpcMetadata.PackageName + "." + sign.ReturnParam + ", error)"),
				Return:      "&" + rpcMetadata.PackageName + "." + sign.ReturnParam + "{} " + ",nil",
			})

			if err != nil {
				t.Fatal(err)
			}
		} else if !sign.IsStreamParam && sign.IsStreamReturnParam {
			result, err = Parser(RPCLogicTemplate, LogicTemplateParam{
				PackageName: rpcMetadata.PackageName,
				MethodName:  sign.Name,
				Project:     "zero",
				Sign:        "(in *" + rpcMetadata.PackageName + "." + sign.Param + ", stream " + rpcMetadata.PackageName + "." + rpcMetadata.ServiceName + "_" + sign.Name + "Server" + ") error",
				Return:      "nil",
			})
			if err != nil {
				t.Fatal(err)
			}
		} else {
			result, err = Parser(RPCLogicTemplate, LogicTemplateParam{
				PackageName: rpcMetadata.PackageName,
				MethodName:  sign.Name,
				Project:     "zero",
				Sign:        "(stream " + rpcMetadata.PackageName + "." + rpcMetadata.ServiceName + "_" + sign.Name + "Server" + ") error",
				Return:      "nil",
			})
			if err != nil {
				t.Fatal(err)
			}
		}
		err = ioutil.WriteFile(fmt.Sprintf("./../../services/%s/rpc/logic/"+"%s_logic.go", rpcMetadata.PackageName, MarshalToSnakeCase(sign.Name)), []byte(result), 0777)
		if err != nil {
			t.Log(err)
		}
		t.Log(result)
	}

}
func TestRPCSVC(t *testing.T) {
	result, err := Parser(RPCSvcTemplate, SvcTemplateParam{
		PackageName: rpcMetadata.PackageName,
		Project:     "zero",
	})

	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("./../../services/hello/rpc/svc/"+"service_context.go", []byte(result), 0777)
	if err != nil {
		t.Log(err)
	}
	t.Log(result)
}
func TestRPCConfig(t *testing.T) {
	err := ioutil.WriteFile("./../../services/hello/rpc/zcfg/"+"zcfg.go", []byte(RPCConfigTemplate), 0777)
	if err != nil {
		t.Log(err)
	}
}
func TestRPCMain(t *testing.T) {
	result, err := Parser(RPCMainTemplate, MainTemplateParam{
		PackageName: rpcMetadata.PackageName,
		Project:     "zero",
	})

	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("./../../services/hello/rpc/"+"main.go", []byte(result), 0777)
	if err != nil {
		t.Log(err)
	}
}
