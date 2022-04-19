package gen

import (
	"github.com/SunMaybo/zero/common/zlog"
	"github.com/SunMaybo/zero/zctl/cmd"
	"github.com/SunMaybo/zero/zctl/file"
	"github.com/SunMaybo/zero/zctl/parser"
	"github.com/SunMaybo/zero/zctl/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func JavaGrpcParentProject(project, groupId, artifactId, version string) {
	path, _ := os.Getwd()
	projectDir := path + "/" + project
	if err := file.MkdirAll(projectDir); err != nil {
		zlog.S.Errorw("create project dir error", "err", err)
		os.Exit(-1)
	}
	parentPom := projectDir + "/pom.xml"
	if ok, err := file.PathExists(parentPom); err != nil {
		zlog.S.Errorw("check parent pom file error", "err", err)
		os.Exit(-1)
	} else if ok {
		zlog.S.Infow("parent pom file exists", "path", parentPom)
	} else {
		if result, err := genParentPom(project, groupId, artifactId, version); err != nil {
			zlog.S.Errorw("generate parent pom file error", "err", err)
			os.Exit(-1)
		} else {
			if err := file.WriterFile(parentPom, []byte(result)); err != nil {
				zlog.S.Errorw("write parent pom file error", "err", err)
				os.Exit(-1)
			}
			zlog.S.Infow("generate parent pom file success", "path", parentPom)
		}
	}
}
func JavaGrpcPackage(project, groupId, artifactId, version string) {
	path, _ := os.Getwd()
	protoProject := path + "/grpc_java/" + project
	if err := file.MkdirAll(file.GetFilePath(protoProject, "/src/main/java")); err != nil {
		zlog.S.Errorw("create proto project dir error", "err", err)
		os.Exit(-1)
	}
	if err := file.MkdirAll(file.GetFilePath(protoProject, "/src/main/resources")); err != nil {
		zlog.S.Errorw("create proto project dir error", "err", err)
		os.Exit(-1)
	}
	if result, err := javaGrpcPom(project, groupId, artifactId, version); err != nil {
		zlog.S.Errorw("generate pom file error", "err", err)
		os.Exit(-1)
	} else {
		if err := file.WriterFile(file.GetFilePath(protoProject, "/pom.xml"), []byte(result)); err != nil {
			zlog.S.Errorw("generate pom file error", "err", err)
			os.Exit(-1)
		}
		zlog.S.Infow("generate pom file success", "path", protoProject+"/pom.xml")
	}
	zlog.S.Infow("generate project success", "path", protoProject)
}
func javaGrpcPom(project, groupId, artifactId, version string) (string, error) {
	result, err := template.Parser(template.ProtoMaven, template.JavaRpcParam{
		GroupId:    groupId,
		ArtifactId: artifactId,
		Project:    project,
		Version:    version,
	})
	if err != nil {
		return "", err
	}
	return result, nil
}
func javaGrpcImpl(packageFile, javaGrpcFilePath string) error {
	serviceSign := parser.ParserJavaGrpc(javaGrpcFilePath)
	if len(serviceSign.MethodSigns) <= 0 {
		return nil
	}
	result, err := template.Parser(template.JavaRPCImplPattern, serviceSign)
	if err != nil {
		return err
	}

	if err := file.WriterFile(file.GetFilePath(packageFile, "/src/main/java/"+strings.ReplaceAll(serviceSign.PackageName, ".", "/")+"/"+serviceSign.ServiceName+".java"), []byte(result)); err != nil {
		zlog.S.Errorw("generate java grpc impl error", "err", err)
		os.Exit(-1)
	}
	return nil
}
func JavaGrpcCompileAndDeploy(mavenBinPath, mavenSettings, protoProjectDir, altDeploymentRepository, workDir string) {
	result, err := cmd.JavaProtoExecute(workDir, protoProjectDir)
	if err != nil {
		zlog.S.Errorw("protoc gen java grpc error", "err", err)
		os.Exit(-1)
	}
	err = filepath.Walk(protoProjectDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		return javaGrpcImpl(protoProjectDir, path)
	})

	if err != nil {
		zlog.S.Errorw("generate java grpc impl error", "err", err)
		os.Exit(-1)
	}
	zlog.S.Infow("protoc gen java grpc success", "result", result)
	if err := JavaGrpcDeploy(mavenBinPath, mavenSettings, protoProjectDir, altDeploymentRepository); err != nil {
		zlog.S.Errorw("mvn clean deploy error", "err", err)
		os.Exit(-1)
	}
	zlog.S.Infow("mvn clean deploy success")
}

func JavaGrpcDeploy(mavenBinPath, mavenSettings, protoProjectDir, altDeploymentRepository string) error {
	return cmd.MavenDeploy(mavenBinPath, mavenSettings, altDeploymentRepository, protoProjectDir)
}
func genParentPom(project, groupId, artifactId, version string) (string, error) {
	result, err := template.Parser(template.ProjectMaven, template.JavaRpcParam{
		GroupId:    groupId,
		ArtifactId: artifactId,
		Project:    project,
		Version:    version,
	})
	if err != nil {
		return "", err
	}
	return result, nil
}
