package gen

import (
	"github.com/SunMaybo/zero/common/zlog"
	"github.com/SunMaybo/zero/zctl/file"
	"github.com/SunMaybo/zero/zctl/template"
	"os"
)

func HttpService(project, module, serviceName string) {
	path, _ := os.Getwd()
	projectDir := path + "/" + "apis" + "/" + module + "/" + serviceName
	if isOk, err := file.PathExists(file.GetFilePath(projectDir, "/svc")); err != nil || !isOk {
		err = file.MkdirAll(file.GetFilePath(projectDir, "/svc"))
		if err != nil {
			zlog.S.Fatalw("create dir error", "err", err)
		}
	}
	if isOk, err := file.PathExists(file.GetFilePath(projectDir, "/svc/server_context.go")); err != nil || !isOk {
		result, err := template.Parser(template.HttpSvcTemplate, &template.HttpSvcTemplateParam{
			Project:     project,
			Module:      module,
			ServiceName: serviceName,
		})
		if err != nil {
			zlog.S.Fatalw("parser template error", "err", err, "")
		}
		_ = file.WriterFile(file.GetFilePath(projectDir, "/svc/server_context.go"), []byte(result))
	}
	if isOk, err := file.PathExists(file.GetFilePath(projectDir, "/server")); err != nil || !isOk {
		err = file.MkdirAll(file.GetFilePath(projectDir, "/server"))
		if err != nil {
			zlog.S.Fatalw("create dir error", "err", err)
		}
	}
	if isOk, err := file.PathExists(file.GetFilePath(projectDir, "/server/server.go")); err != nil || !isOk {
		result, err := template.Parser(template.HttpServerTemplate, &template.HttpServerTemplateParam{
			Project:     project,
			Module:      module,
			ServiceName: serviceName,
		})
		if err != nil {
			zlog.S.Fatalw("parser template error", "err", err, "")
		}
		_ = file.WriterFile(file.GetFilePath(projectDir, "/server/server.go"), []byte(result))
	}
	if isOk, err := file.PathExists(file.GetFilePath(projectDir, "/etc")); err != nil || !isOk {
		err = file.MkdirAll(file.GetFilePath(projectDir, "/etc"))
		if err != nil {
			zlog.S.Fatalw("create dir error", "err", err)
		}
	}
	if isOk, err := file.PathExists(file.GetFilePath(projectDir, "/etc/config.yaml")); err != nil || !isOk {
		_ = file.WriterFile(file.GetFilePath(projectDir, "/etc/config.yaml"), []byte(template.HttpETCTemplate))
	}

	if isOk, err := file.PathExists(file.GetFilePath(projectDir, "/config")); err != nil || !isOk {
		err = file.MkdirAll(file.GetFilePath(projectDir, "/config"))
		if err != nil {
			zlog.S.Fatalw("create dir error", "err", err)
		}
	}
	if isOk, err := file.PathExists(file.GetFilePath(projectDir, "/config/config.go")); err != nil || !isOk {
		_ = file.WriterFile(file.GetFilePath(projectDir, "/config/config.go"), []byte(template.HttpConfigTemplate))
	}
	if isOk, err := file.PathExists(file.GetFilePath(projectDir, "/main.go")); err != nil || !isOk {
		result, err := template.Parser(template.HttpMainTemplate, &template.HttpMainTemplateParam{
			Project:     project,
			Module:      module,
			ServiceName: serviceName,
		})
		if err != nil {
			zlog.S.Fatalw("parser template error", "err", err, "")
		}
		_ = file.WriterFile(file.GetFilePath(projectDir, "/main.go"), []byte(result))
	}
	return
}
