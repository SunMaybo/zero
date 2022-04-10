package gen

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/SunMaybo/zero/zctl/parser"
	"github.com/SunMaybo/zero/zctl/template"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const serverPathPattern = "/%s/rpc/server/" + "server.go"
const logicPathPattern = "/%s/rpc/logic/" + "%s_logic.go"
const svcPathPattern = "/%s/rpc/svc/service_context.go"
const rpcConfigPathPattern = "/%s/rpc/config/config.go"
const rpcEtcPathPattern = "/%s/rpc/etc/config.yaml"
const rpcMainPathPattern = "/%s/rpc/" + "main.go"

type RpcBuilder struct {
	project     string
	projectPath string
	module      string
	serviceType string
}

func NewRpcBuilder(project, projectPath, module, serviceType string) *RpcBuilder {
	return &RpcBuilder{
		project:     project,
		module:      module,
		projectPath: projectPath,
		serviceType: serviceType,
	}
}
func (r *RpcBuilder) StartBuild() {
	filepath.Walk(r.projectPath+"/proto/"+r.serviceType+"/"+r.module, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), "proto") {
			//protoc  --go_out=proto/test_services  proto/$(services)/greeter.proto
			_, err := Run(fmt.Sprintf("protoc -I=%s %s --go_out=%s --go-grpc_out=%s", r.projectPath+"/proto/", path, r.projectPath+"/proto/"+r.serviceType+"/"+r.module, r.projectPath+"/proto/"+r.serviceType+"/"+r.module), r.projectPath+"/proto/"+r.serviceType+"/"+r.module)
			if err != nil {
				panic(err)
			}
			if rpcMetadata, err := parser.Parser(path); err != nil {
				panic(err)
			} else {
				os.MkdirAll(r.projectPath+"/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName+"/rpc", 0777)
				os.MkdirAll(r.projectPath+"/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName+"/rpc/config", 0777)
				os.MkdirAll(r.projectPath+"/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName+"/rpc/etc", 0777)
				os.MkdirAll(r.projectPath+"/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName+"/rpc/logic", 0777)
				os.MkdirAll(r.projectPath+"/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName+"/rpc/server", 0777)
				os.MkdirAll(r.projectPath+"/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName+"/rpc/svc", 0777)
				os.MkdirAll(r.projectPath+"/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName+"/rpc/svc", 0777)
				os.MkdirAll(r.projectPath+"/docs/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName, 0777)
				_, err = Run(fmt.Sprintf("protoc-go-inject-tag -input=%s", r.projectPath+"/proto/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName+"/"+rpcMetadata.PackageName+".pb.go"), r.projectPath+"/proto/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName)
				if err != nil {
					panic(err)
				}
				_, err = Run(fmt.Sprintf("protoc --doc_out=%s --doc_opt=html,index.html *.proto", r.projectPath+"/docs/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName), r.projectPath+"/proto/"+r.serviceType+"/"+r.module)
				if err != nil {
					panic(err)
				}
				_, err = Run(fmt.Sprintf("protoc --doc_out=%s --doc_opt=markdown,index.md *.proto", r.projectPath+"/docs/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName), r.projectPath+"/proto/"+r.serviceType+"/"+r.module)
				if err != nil {
					panic(err)
				}
				if err := r.rpcMethod(rpcMetadata.PackageName, rpcMetadata); err != nil {
					panic(err)
				}
				if err := r.rpcSVC(r.project, rpcMetadata.PackageName); err != nil {
					panic(err)
				}
				if err := r.rpcConfig(rpcMetadata.PackageName); err != nil {
					panic(err)
				}
				if err := r.rpcServer(r.project, rpcMetadata); err != nil {
					panic(err)
				}
				if err := r.rpcLogic(rpcMetadata); err != nil {
					panic(err)
				}
				if err := r.rpcMain(r.project, rpcMetadata.PackageName, rpcMetadata.ServiceName); err != nil {
					panic(err)
				}
			}
		}
		return nil
	})
}

const (
	// OsWindows represents os windows
	OsWindows = "windows"
	// OsMac represents os mac
	OsMac = "darwin"
	// OsLinux represents os linux
	OsLinux = "linux"
	// OsJs represents os js
	OsJs = "js"
	// OsIOS represents os ios
	OsIOS = "ios"
)

func Run(arg, dir string) (string, error) {
	goos := runtime.GOOS
	var cmd *exec.Cmd
	switch goos {
	case OsMac, OsLinux:
		cmd = exec.Command("/bin/bash", "-c", arg)
	case OsWindows:
		cmd = exec.Command("cmd.exe", "/c", arg)
	default:
		return "", fmt.Errorf("unexpected os: %v", goos)
	}
	if len(dir) > 0 {
		cmd.Dir = dir
	}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return "", errors.New(stderr.String())
		}
		return "", err
	}

	return stdout.String(), nil
}

func (r *RpcBuilder) rpcServer(project string, rpcMetadata parser.RpcMetadata) error {
	temps := template.ServerTemplateParam{
		Project:     project,
		Module:      r.module,
		ServiceType: r.serviceType,
		ServiceName: rpcMetadata.ServiceName,
		PackageName: rpcMetadata.PackageName,
	}
	for _, sign := range rpcMetadata.MethodSigns {
		ms := template.MethodSign{}
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
	result, err := template.Parser(template.RPCServerTemplate, temps)
	if err != nil {
		return err
	}
	serverPath := r.projectPath + "/" + r.serviceType + "/" + r.module + "/" + fmt.Sprintf(serverPathPattern, temps.PackageName)
	err = ioutil.WriteFile(serverPath, []byte(result), 0777)
	if err != nil {
		return err
	}
	return nil
}
func (r *RpcBuilder) rpcLogic(rpcMetadata parser.RpcMetadata) error {
	for _, sign := range rpcMetadata.MethodSigns {
		var result string
		var err error
		if !sign.IsStreamReturnParam && !sign.IsStreamParam {
			result, err = template.Parser(template.RPCLogicTemplate, template.LogicTemplateParam{
				PackageName: rpcMetadata.PackageName,
				MethodName:  sign.Name,
				Module:      r.module,
				ServiceType: r.serviceType,
				Project:     r.project,
				Sign:        ("(in *" + rpcMetadata.PackageName + "." + sign.Param + ")") + (" (*" + rpcMetadata.PackageName + "." + sign.ReturnParam + ", error)"),
				Return:      "&" + rpcMetadata.PackageName + "." + sign.ReturnParam + "{} " + ",nil",
			})

			if err != nil {
				return err
			}
		} else if !sign.IsStreamParam && sign.IsStreamReturnParam {
			result, err = template.Parser(template.RPCLogicTemplate, template.LogicTemplateParam{
				PackageName: rpcMetadata.PackageName,
				MethodName:  sign.Name,
				ServiceType: r.serviceType,
				Project:     r.project,
				Sign:        "(in *" + rpcMetadata.PackageName + "." + sign.Param + ", stream " + rpcMetadata.PackageName + "." + rpcMetadata.ServiceName + "_" + sign.Name + "Server" + ") error",
				Return:      "nil",
				Module:      r.module,
			})
			if err != nil {
				return err
			}
		} else {
			result, err = template.Parser(template.RPCLogicTemplate, template.LogicTemplateParam{
				PackageName: rpcMetadata.PackageName,
				MethodName:  sign.Name,
				ServiceType: r.serviceType,
				Project:     r.project,
				Sign:        "(stream " + rpcMetadata.PackageName + "." + rpcMetadata.ServiceName + "_" + sign.Name + "Server" + ") error",
				Return:      "nil",
				Module:      r.module,
			})
			if err != nil {
				return err
			}
		}
		path := r.projectPath + "/" + r.serviceType + "/" + r.module + "/" + fmt.Sprintf(logicPathPattern, rpcMetadata.PackageName, template.MarshalToSnakeCase(sign.Name))
		if ok, err := pathExists(path); err != nil {
			return err
		} else if !ok {
			err := ioutil.WriteFile(path, []byte(result), 0777)
			if err != nil {
				return err
			}
		}
	}
	return nil

}
func (r *RpcBuilder) rpcMethod(packageName string, rpcMetadata parser.RpcMetadata) error {
	names := template.MethodName{
		PackageName: packageName,
	}
	for _, sign := range rpcMetadata.MethodSigns {
		names.Names = append(names.Names, sign.Name)
	}
	result, err := template.Parser(template.RPCMethodTemplate, names)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(r.projectPath+"/proto/"+r.serviceType+"/"+r.module+"/"+packageName+"/method.go", []byte(result), 0777)
}
func (r *RpcBuilder) rpcSVC(project, packageName string) error {
	result, err := template.Parser(template.RPCSvcTemplate, template.SvcTemplateParam{
		PackageName: packageName,
		Project:     project,
		Module:      r.module,
		ServiceType: r.serviceType,
	})
	if err != nil {
		return err
	}
	if ok, err := pathExists(r.projectPath + "/" + r.serviceType + "/" + r.module + "/" + fmt.Sprintf(svcPathPattern, packageName)); err != nil {
		return err
	} else if !ok {
		return ioutil.WriteFile(r.projectPath+"/"+r.serviceType+"/"+r.module+"/"+fmt.Sprintf(svcPathPattern, packageName), []byte(result), 0777)
	}
	return nil
}
func (r *RpcBuilder) rpcConfig(packageName string) error {
	configfile := r.projectPath + "/" + r.serviceType + "/" + r.module + "/" + fmt.Sprintf(rpcConfigPathPattern, packageName)
	if ok, err := pathExists(configfile); err != nil {
		return err
	} else if !ok {
		ioutil.WriteFile(configfile, []byte(template.RPCConfigTemplate), 0777)
		result, err := template.Parser(template.RPCETCTemplate, template.EtcConfig{
			PackageName: packageName,
		})
		if err != nil {
			return err
		}
		ioutil.WriteFile(r.projectPath+"/"+r.serviceType+"/"+r.module+"/"+fmt.Sprintf(rpcEtcPathPattern, packageName), []byte(result), 0777)
	}
	return nil
}
func (r *RpcBuilder) rpcMain(project, packageName, service string) error {
	result, err := template.Parser(template.RPCMainTemplate, template.MainTemplateParam{
		PackageName: packageName,
		Project:     project,
		Module:      r.module,
		ServiceType: r.serviceType,
		Service:     template.MarshalToCamelCase(service),
	})

	if err != nil {
		return err
	}
	if ok, err := pathExists(r.projectPath + "/" + r.serviceType + "/" + r.module + "/" + fmt.Sprintf(rpcMainPathPattern, packageName)); err != nil {
		return err
	} else if !ok {
		return ioutil.WriteFile(r.projectPath+"/"+r.serviceType+"/"+r.module+"/"+fmt.Sprintf(rpcMainPathPattern, packageName), []byte(result), 0777)
	}
	return nil
}
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
