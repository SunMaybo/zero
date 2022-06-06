package execute

import (
	"errors"
	"github.com/SunMaybo/zero/common/zlog"
	"github.com/SunMaybo/zero/zctl/config"
	"github.com/SunMaybo/zero/zctl/file"
	"github.com/SunMaybo/zero/zctl/gen"
	"github.com/SunMaybo/zero/zctl/install"
	"github.com/SunMaybo/zero/zctl/parser"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	groupId                 *string
	artifactId              *string
	versionGrpc             *string
	protoDir                *string
	altDeploymentRepository *string
	maven                   *string
	proxy                   *string
	installLang             *string
	golangModule            *string
	golangServiceType       *string
	mavenSettings           *string
	docSource               *string
	docType                 *int
	moduleHttp              *string
	serviceHttp             *string
)

var genProjectCommand = &cobra.Command{
	Use:   "java_project",
	Short: "generate Maven Parent project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gen.JavaGrpcParentProject(*artifactId, *groupId, *artifactId, "0.0.1-SNAPSHOT")
	},
}
var genDocCommand = &cobra.Command{
	Use:   "doc",
	Short: "generate proto doc",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		baseDir, _ := os.Getwd()
		docDst := file.GetFilePath(baseDir, "/docs")
		if exist, err := file.PathExists(docDst); err != nil || !exist {
			err = file.MkdirAll(docDst)
			if err != nil {
				panic(err)
			}
		}
		sourceRelative := ""
		if !strings.HasPrefix(*docSource, "/") {
			sourceRelative = file.GetFilePath(baseDir, "/"+*docSource)
		}
		_ = file.MkdirAll(file.GetFilePath(docDst, "/"+*docSource))
		gen.GenDoc(sourceRelative, file.GetFilePath(docDst, "/"+*docSource), *docType)
	},
}

var genGrpcCommand = &cobra.Command{
	Use:   "java_grpc_package",
	Short: "generate Maven grpc package",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		workDir := ""
		path, _ := os.Getwd()
		if strings.HasPrefix(*protoDir, "/") {
			workDir = file.GetFilePath(path, *protoDir)
		} else {
			workDir = file.GetFilePath(path, "/"+*protoDir)
		}
		//通过proto文件获取当前java package
		javaPackage := ""
		if exist, err := file.PathExists(workDir); err != nil || !exist {
			zlog.S.Infow("proto dir not exist", "path", workDir)
			return
		}
		err := filepath.Walk(workDir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(path) == ".proto" {
				protoFileMd, _ := parser.Parser(path)
				if javaPackage != "" && javaPackage != protoFileMd.JavaPackageName {
					return errors.New("proto file java_package is not same")
				}
				javaPackage = protoFileMd.JavaPackageName
			}
			return nil
		})
		if err != nil {
			zlog.S.Error(err)
			return
		}
		if javaPackage == "" {
			zlog.S.Error("can not find java package")
			return
		}
		if !strings.Contains(javaPackage, ".proto.") {
			zlog.S.Errorf("java package is not correct,current package is %s", javaPackage)
			return
		}
		groupIdGrpc := javaPackage[:strings.Index(javaPackage, ".proto.")]
		project := javaPackage[strings.Index(javaPackage, ".proto.")+7:]
		if len(project) <= 0 {
			zlog.S.Error("can not find java package")
			return
		}
		protoProject := file.GetFilePath(path, "/grpc_java/"+project)
		artifactIdGrpc := strings.ReplaceAll(project, "_", "-")
		gen.JavaGrpcPackage(project, groupIdGrpc, artifactIdGrpc, *versionGrpc)
		gen.JavaGrpcCompileAndDeploy(*maven, *mavenSettings, protoProject, *altDeploymentRepository, workDir)
		zlog.S.Infow("generate grpc package success", "groupId", groupIdGrpc, "artifactId", artifactIdGrpc, "version", *versionGrpc)
	},
}
var installCommand = &cobra.Command{
	Use:   "install",
	Short: "install protoc、grpc、validate、doc",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if *proxy != "" {
			_ = os.Setenv("http_proxy", *proxy)
			_ = os.Setenv("https_proxy", *proxy)
		}
		installLang := *installLang
		if installLang == "java" {
			install.JavaInstall()
		} else if installLang == "golang" {
			install.GolangInstall()
		} else {
			zlog.S.Errorf("install language %s not support", installLang)
		}
		zlog.S.Infow("install success", "language", installLang)
	},
}
var golangModuleCommand = &cobra.Command{
	Use:   "golang_module",
	Short: "generate golang module",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		workDir, _ := os.Getwd()
		b := gen.NewRpcBuilder(getGolangProjectByMod(), workDir, *golangModule, *golangServiceType)
		b.StartBuild()
	},
}
var golangHttpModuleCommand = &cobra.Command{
	Use:   "golang_http_service",
	Short: "generate golang http service",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		gen.HttpService(getGolangProjectByMod(), *moduleHttp, *serviceHttp)
	},
}

func init() {
	groupId = genProjectCommand.Flags().String("g", "", "Maven:groupId")
	artifactId = genProjectCommand.Flags().String("a", "", "Maven:artifactId")

	versionGrpc = genGrpcCommand.Flags().String("v", "0.0.1-SNAPSHOT", "Maven:version")
	protoDir = genGrpcCommand.Flags().String("p", "", "proto dir")
	altDeploymentRepository = genGrpcCommand.Flags().String("r",
		"",
		"alt deployment repository")

	proxy = installCommand.Flags().String("proxy", "", "proxy")
	installLang = installCommand.Flags().String("lang", "java", "install language")

	golangModule = golangModuleCommand.Flags().String("m", "greeter", "golang module")
	golangServiceType = golangModuleCommand.Flags().String("t", "services", "golang service type")

	maven = genGrpcCommand.Flags().String("m", "", "maven exec path")

	docSource = genDocCommand.Flags().String("s", "", "doc source")
	docType = genDocCommand.Flags().Int("t", 1, "doc type")

	moduleHttp = golangHttpModuleCommand.Flags().String("m", "", "golang module")
	serviceHttp = golangHttpModuleCommand.Flags().String("s", "", "golang service")
}
func GetAllCommands(cfg config.Config) []*cobra.Command {
	if *maven == "" {
		maven = &cfg.Maven
	}
	if *altDeploymentRepository == "" {
		altDeploymentRepository = &cfg.MavenDeploymentRepository

	}
	if *proxy == "" && cfg.Proxy != "" {
		proxy = &cfg.Proxy
	}
	mavenSettings = &cfg.MavenSettings
	return []*cobra.Command{
		genProjectCommand,
		genGrpcCommand,
		installCommand,
		golangModuleCommand,
		genDocCommand,
		golangHttpModuleCommand,
	}
}
func getGolangProjectByMod() string {
	buff, err := ioutil.ReadFile("go.mod")
	if err != nil {
		zlog.S.Errorf("read go.mod error:%v", err)
		os.Exit(-1)
	}
	str := string(buff)
	for _, s := range strings.Split(str, "\n") {
		if strings.HasPrefix(strings.TrimSpace(s), "module") {
			return strings.TrimSpace(strings.Split(s, " ")[1])
		}
	}
	zlog.S.Error("can not find module from go.mod")
	os.Exit(-1)
	return ""
}
