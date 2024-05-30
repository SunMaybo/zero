package gen

import (
	"bufio"
	"fmt"
	"github.com/SunMaybo/zero/common/zlog"
	"github.com/SunMaybo/zero/zctl/cmd"
	"github.com/SunMaybo/zero/zctl/file"
	"github.com/SunMaybo/zero/zctl/install"
	"github.com/SunMaybo/zero/zctl/parser"
	"github.com/SunMaybo/zero/zctl/template"
	"io"
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

	if err := file.MkdirAll(file.GetFilePath(protoProject, "/src/main/java/cn/zero/grpc/proto/extend")); err != nil {
		zlog.S.Errorw("create proto extend dir error", "err", err)
		os.Exit(-1)
	}
	if err := file.WriterFile(file.GetFilePath(protoProject, "/src/main/java/cn/zero/grpc/proto/extend/Extend.java"), []byte(install.EXTEND_GRPC_CLAZZ)); err != nil {
		zlog.S.Errorw("create Extend.java error", "err", err)
		os.Exit(-1)
	}
	if err := file.WriterFile(file.GetFilePath(protoProject, "/src/main/java/cn/zero/grpc/proto/extend/ExtendValidator.java"), []byte(install.EXTEND_GRPC_VALIDATE_CLAZZ)); err != nil {
		zlog.S.Errorw("create ExtendValidator.java error", "err", err)
		os.Exit(-1)
	}

	if err := file.MkdirAll(file.GetFilePath(protoProject, "/src/main/java/cn/zero/grpc/proto/xbb")); err != nil {
		zlog.S.Errorw("create proto extend dir error", "err", err)
		os.Exit(-1)
	}
	if err := file.WriterFile(file.GetFilePath(protoProject, "/src/main/java/cn/zero/grpc/proto/xbb/Xbb.java"), []byte(install.XBB_GRPC_CLAZZ)); err != nil {
		zlog.S.Errorw("create Extend.java error", "err", err)
		os.Exit(-1)
	}
	if err := file.WriterFile(file.GetFilePath(protoProject, "/src/main/java/cn/zero/grpc/proto/xbb/XbbValidator.java"), []byte(install.XBB_VALIDATE_CLAZZ)); err != nil {
		zlog.S.Errorw("create ExtendValidator.java error", "err", err)
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

	filepath.Walk(workDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".proto") {
			zlog.S.Infow("cp .proto:" + path + "," + info.Name())
			CopyFile(protoProjectDir+"/src/main/resources/"+info.Name(), path)
		}
		return nil
	})
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
		SwitchGrpcJavaType(path)
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

var replaceMap = map[string]string{
	"(double value)": "(java.lang.Double value)",
	"(int value)":    "(java.lang.Integer value)",
	"(long value)":   "(java.lang.Long value)",
	"(float value)":  "(java.lang.Float value)",
}

var appendFuncMap = map[string]string{
	"(double value)": "public Builder %sValue(java.lang.Double value) {if (value==null) return this;return %s(value);}",
	"(int value)":    "public Builder %sValue(java.lang.Integer value) {if (value==null) return this;return %s(value);}",
	"(long value)":   "public Builder %sValue(java.lang.Long value) {if (value==null) return this;return %s(value);}",
	"(float value)":  "public Builder %sValue(java.lang.Float value) {if (value==null) return this;return %s(value);}",
	"String":         "public Builder %sValue(java.lang.String value) {if (value==null) return this;return %s(value);}",
}

var filterBySuffix = map[string]struct{}{
	"Builder":   {},
	"Validator": {},
	"Grpc":      {},
	"Proto":     {},
}

const prefix = "public Builder set"

func SwitchGrpcJavaType(path string) {
	for s := range filterBySuffix {
		if strings.HasSuffix(path, s) {
			return
		}
	}

	buff, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var appendFunc []string
	items := strings.Split(string(buff), "\n")
	var result []string
	for i := 0; i < len(items); i++ {
		item := items[i]
		if item == "" {
			result = append(result, item)
		} else if strings.HasPrefix(strings.TrimSpace(item), "java.lang.String value) {") {
			if len(items) > i+2 && strings.HasPrefix(strings.TrimSpace(items[i+1]), "if (value == null) {") && strings.Contains(items[i+2], "throw new NullPointerException();") {
				result = append(result, items[i])
				result = append(result, items[i+1])
				if name, isOk := parseFuncName(item); isOk {
					appendFunc = append(appendFunc, fmt.Sprintf(appendFuncMap["String"], name, name))
				}
				i += 2
				continue
			} else {
				result = append(result, item)
				continue
			}
		} else if !strings.HasPrefix(strings.TrimSpace(item), prefix) {
			result = append(result, item)
			continue
		}
		isAppend := false
		for k := range replaceMap {
			if strings.Contains(item, k) {
				if name, isOk := parseFuncName(item); isOk {
					appendFunc = append(appendFunc, fmt.Sprintf(appendFuncMap[k], name, name))
				}
				result = append(result, item)
				isAppend = true
				break
			}
		}
		if !isAppend {
			result = append(result, item)
		}
	}
	result = append(result, appendFunc...)
	_ = os.WriteFile(path, []byte(strings.Join(result, "\n")), 0777)
}
func parseFuncName(item string) (string, bool) {
	if strings.Contains(item, prefix) {
		beginIdx := strings.Index(item, prefix)
		if beginIdx <= 0 {
			return "", false
		}
		item = item[beginIdx:]
		endIdx := strings.Index(item, "(")
		if endIdx <= 0 {
			return "", false
		}
		return "set" + item[:endIdx], true
	}
	return "", false
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
func CopyFile(dstFilePath string, srcFilePath string) (written int64, err error) {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		fmt.Printf("打开源文件错误，错误信息=%v\n", err)
	}
	defer srcFile.Close()
	reader := bufio.NewReader(srcFile)

	dstFile, err := os.OpenFile(dstFilePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Printf("打开目标文件错误，错误信息=%v\n", err)
		return
	}
	writer := bufio.NewWriter(dstFile)
	defer dstFile.Close()
	return io.Copy(writer, reader)
}
