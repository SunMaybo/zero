package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/SunMaybo/zero/common/zlog"
	"github.com/SunMaybo/zero/zctl/file"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

const (
	// OsWindows represents os windows
	OsWindows = "windows"
	// OsMac represents os mac
	OsMac = "darwin"
	// OsLinux represents os linux
	OsLinux = "linux"
)

type DocType int

const (
	Html = iota
	Markdown
)

func Run(arg, dir string) (string, error) {
	zlog.S.Infof("run command: %s,dir:%s", arg, dir)
	goos := runtime.GOOS
	var cmd *exec.Cmd
	switch goos {
	case OsMac, OsLinux:
		cmd = exec.Command("/bin/bash", "-c", arg)
	case OsWindows:
		cmd = exec.Command(arg)
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
func JavaProtoExecute(workDir, protoProjectDir string) (string, error) {
	if runtime.GOOS == "windows" {
		workDir = strings.ReplaceAll(workDir, "\\.\\", "\\")
		workDir = strings.ReplaceAll(workDir, "/./", "/")
		return windowExec(getProtoc(), workDir,
			"--experimental_allow_proto3_optional",
			fmt.Sprintf("--plugin=protoc-gen-grpc-java=%s", getProtoJavaGrpc()),
			fmt.Sprintf("--plugin=protoc-gen-validate=%s", getProtoJavaValidate()),
			fmt.Sprintf("-I=%s", getProtoInclude()),
			fmt.Sprintf("-I=%s", workDir),
			fmt.Sprintf("--validate_out=lang=java:%s", file.GetFilePath(protoProjectDir, "/src/main/java")),
			fmt.Sprintf("--java_out=%s", file.GetFilePath(protoProjectDir, "/src/main/java")),
			fmt.Sprintf("--grpc-java_out=%s", file.GetFilePath(protoProjectDir, "/src/main/java")),
			"*.proto",
		)
	} else {
		return Run(getProtoc()+fmt.Sprintf(
			" --experimental_allow_proto3_optional --plugin=protoc-gen-grpc-java=%s "+
				"--plugin=protoc-gen-validate=%s "+
				"-I=%s "+
				"-I=%s "+
				"--validate_out=lang=java:%s --java_out=%s  --grpc-java_out=%s *.proto",
			getProtoJavaGrpc(), getProtoJavaValidate(), getProtoInclude(), workDir,
			file.GetFilePath(protoProjectDir, "/src/main/java"),
			file.GetFilePath(protoProjectDir, "/src/main/java"),
			file.GetFilePath(protoProjectDir, "/src/main/java")), workDir)
	}

}
func windowExec(protoc string, dir string, args ...string) (string, error) {
	zlog.S.Infof("protoc:%s,%s", dir, protoc+" "+strings.Join(args, " "))
	cmd := exec.Command(protoc, args...)
	cmd.Dir = dir
	result := ""
	if outPipe, err := cmd.StdoutPipe(); err != nil {
		return "", err
	} else {
		cmd.Start()
		in := bufio.NewScanner(outPipe)
		for in.Scan() {
			var buff []byte
			for i, b := range in.Bytes() {
				if b == 0 {
					continue
				}
				if len(in.Bytes()) <= i+2 && b == 13 {
					continue
				}
				buff = append(buff, b)
			}
			result += string(buff)
		}
		cmd.Wait()
	}
	return result, nil
}
func GolangProtoExecute(workDir, protoProjectDir, protoFilePath string) (string, error) {
	return Run(getProtoc()+fmt.Sprintf(" --plugin=protoc-gen-go-grpc=%s -I=%s -I=%s --go_out=%s --go-grpc_out=%s %s",
		getProtoGolangGrpc(),
		getProtoInclude(),
		protoProjectDir,
		protoProjectDir,
		protoProjectDir,
		protoFilePath), workDir)
}
func GetGolangProtoValidate(protoServicePath, pbFilePath string) error {
	_, err := Run(getProtoGolangInject()+fmt.Sprintf(" -input=%s", pbFilePath), protoServicePath)
	return err
}

func GetProtoDoc(workDir, docOutDir string, docType DocType) error {
	if docType == Html {
		_, err := Run(getProtoc()+fmt.Sprintf(" --doc_out=%s --plugin=protoc-gen-doc=%s --doc_opt=html,index.html *.proto", docOutDir, getProtoDoc()), workDir)
		return err
	} else {
		_, err := Run(getProtoc()+fmt.Sprintf(" --doc_out=%s --plugin=protoc-gen-doc=%s --doc_opt=markdown,index.md *.proto", docOutDir, getProtoDoc()), workDir)
		return err
	}

}
func getProtoc() string {
	return file.GetFilePath(getProtoDir(), "/bin/protoc")
}
func getProtoInclude() string {
	return file.GetFilePath(getProtoDir(), "/include")
}
func getProtoDoc() string {
	return file.GetFilePath(getProtoDir(), "/bin/protoc-gen-doc")
}
func getProtoJavaGrpc() string {
	return file.GetFilePath(getProtoDir(), "/bin/protoc-gen-java-grpc")
}
func getProtoJavaValidate() string {
	return file.GetFilePath(getProtoDir(), "/bin/protoc-gen-validate")
}
func getProtoGolangInject() string {
	golangBinPath, err := GetGolangBinPath()
	if err != nil {
		zlog.S.Errorf("get golang bin path error:%v", err)
		os.Exit(1)
	}
	return file.GetFilePath(golangBinPath, "/protoc-go-inject-tag")
}
func getProtoGolangGrpc() string {
	golangBinPath, err := GetGolangBinPath()
	if err != nil {
		zlog.S.Errorf("get golang bin path error:%v", err)
		os.Exit(1)
	}
	return file.GetFilePath(golangBinPath, "/protoc-gen-go-grpc")
}
func getProtoDir() string {
	path := "/usr/local"
	u, _ := user.Current()
	if runtime.GOOS == "windows" {
		path = u.HomeDir + "\\proto"
	} else {
		path = u.HomeDir + "/proto"
	}
	return path
}
