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
	"runtime"
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
func javaGrpcImpl(javaGrpcDir, javaGrpcFileName string) error {
	serviceSign := parser.ParserJavaGrpc(file.GetFilePath(javaGrpcDir, javaGrpcFileName))
	if len(serviceSign.MethodSigns) <= 0 {
		return nil
	}
	packageName := strings.ReplaceAll(javaGrpcDir, "//", "/")
	packageName = strings.ReplaceAll(packageName, "\\", "/")
	packageName = strings.ReplaceAll(packageName, "\\\\", "/")
	packageName = strings.ReplaceAll(packageName, "/", ".")
	packageName = packageName[strings.Index(packageName, "src.main.java.")+len("src.main.java."):]
	if strings.HasSuffix(packageName, ".") {
		packageName = packageName[:len(packageName)-1]
	}
	serviceSign.PackageName = packageName
	result, err := template.Parser(template.JavaRPCImplPattern, serviceSign)
	if err != nil {
		return err
	}

	if err := file.WriterFile(file.GetFilePath(javaGrpcDir, serviceSign.ServiceName+".java"), []byte(result)); err != nil {
		zlog.S.Errorw("generate java grpc impl error", "err", err)
		os.Exit(-1)
	}
	return nil
}
func JavaGrpcCompileAndDeploy(mavenBinPath, mavenSettings, protoProjectDir, altDeploymentRepository, workDir string) {
	_, err := cmd.Run("cp -rf "+workDir+" "+protoProjectDir+"/src/main/resources", workDir)
	if err != nil {
		zlog.S.Errorw("cp *.proto err", err)
		os.Exit(-1)
	}
	result, err := cmd.JavaProtoExecute(workDir, protoProjectDir)
	if err != nil {
		zlog.S.Errorw("protoc gen java grpc error", "err", err)
		os.Exit(-1)
	} else {
		zlog.S.Infow("proto gen java grpc result", "result", result)
	}
	err = filepath.Walk(protoProjectDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if runtime.GOOS == "windows" {
			path = strings.ReplaceAll(path, "\\", "/")
		}

		return javaGrpcImpl(strings.ReplaceAll(path, info.Name(), ""), info.Name())
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
