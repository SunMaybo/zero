package gen

import (
	"fmt"
	"github.com/SunMaybo/zero/zctl/cmd"
	"github.com/SunMaybo/zero/zctl/file"
	"github.com/SunMaybo/zero/zctl/parser"
	"github.com/SunMaybo/zero/zctl/template"
	"io/fs"
	"path/filepath"
	"strings"
)

const (
	serverPathPattern    = "/%s/rpc/server/" + "server.go"
	logicPathPattern     = "/%s/rpc/logic/" + "%s_logic.go"
	svcPathPattern       = "/%s/rpc/svc/service_context.go"
	rpcConfigPathPattern = "/%s/rpc/config/config.go"
	rpcEtcPathPattern    = "/%s/rpc/etc/config.yaml"
	rpcMainPathPattern   = "/%s/rpc/" + "main.go"
	rpcMethodPathPattern = "/%s/method.go"
)

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
	filepath.Walk(file.GetFilePath(r.projectPath, "/proto/"+r.serviceType+"/"+r.module), func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), "proto") {
			workDir := file.GetFilePath(r.projectPath, "/proto/"+r.serviceType+"/"+r.module)
			if _, err = cmd.GolangProtoExecute(workDir, workDir, path); err != nil {
				panic(err)
			}
			if rpcMetadata, err := parser.Parser(path); err != nil {
				panic(err)
			} else {
				_ = file.MkdirAll(file.GetFilePath(r.projectPath, "/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName+"/rpc/config"))
				_ = file.MkdirAll(file.GetFilePath(r.projectPath, "/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName+"/rpc/etc"))
				_ = file.MkdirAll(file.GetFilePath(r.projectPath, "/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName+"/rpc/logic"))
				_ = file.MkdirAll(file.GetFilePath(r.projectPath, "/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName+"/rpc/server"))
				_ = file.MkdirAll(file.GetFilePath(r.projectPath, "/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName+"/rpc/svc"))
				err = cmd.GetGolangProtoValidate(file.GetFilePath(r.projectPath, "/proto/"+r.serviceType+"/"+r.module+"/"+rpcMetadata.PackageName),
					file.GetFilePath(workDir, "/"+rpcMetadata.PackageName+"/"+strings.ReplaceAll(info.Name(), ".proto", "")+".pb.go"))
				if err != nil {
					panic(err)
				}
				if rpcMetadata.ServiceName == "" {
					return nil
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
	serverPath := file.GetFilePath(r.getModulePath(), fmt.Sprintf(serverPathPattern, temps.PackageName))
	err = file.WriterFile(serverPath, []byte(result))
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
		path := file.GetFilePath(r.getModulePath(), fmt.Sprintf(logicPathPattern, rpcMetadata.PackageName, template.MarshalToSnakeCase(sign.Name)))
		if ok, err := file.PathExists(path); err != nil {
			return err
		} else if !ok {
			err := file.WriterFile(path, []byte(result))
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
	return file.WriterFile(file.GetFilePath(r.projectPath, "/proto/"+r.serviceType+"/"+r.module+fmt.Sprintf(rpcMethodPathPattern, rpcMetadata.PackageName)), []byte(result))
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
	if ok, err := file.PathExists(file.GetFilePath(r.getModulePath(), fmt.Sprintf(svcPathPattern, packageName))); err != nil {
		return err
	} else if !ok {
		return file.WriterFile(file.GetFilePath(r.getModulePath(), fmt.Sprintf(svcPathPattern, packageName)), []byte(result))
	}
	return nil
}
func (r *RpcBuilder) rpcConfig(packageName string) error {
	configFile := file.GetFilePath(r.getModulePath(), fmt.Sprintf(rpcConfigPathPattern, packageName))
	if ok, err := file.PathExists(configFile); err != nil {
		return err
	} else if !ok {
		_ = file.WriterFile(configFile, []byte(template.RPCConfigTemplate))
		result, err := template.Parser(template.RPCETCTemplate, template.EtcConfig{
			PackageName: packageName,
		})
		if err != nil {
			return err
		}
		_ = file.WriterFile(file.GetFilePath(r.getModulePath(), fmt.Sprintf(rpcEtcPathPattern, packageName)), []byte(result))
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
	servicePath := file.GetFilePath(r.getModulePath(), fmt.Sprintf(rpcMainPathPattern, packageName))
	if ok, err := file.PathExists(servicePath); err != nil {
		return err
	} else if !ok {
		return file.WriterFile(servicePath, []byte(result))
	}
	return nil
}
func (r *RpcBuilder) getModulePath() string {
	return file.GetFilePath(r.projectPath, "/"+r.serviceType+"/"+r.module)

}
