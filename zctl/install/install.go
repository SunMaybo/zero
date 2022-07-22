package install

import (
	"github.com/SunMaybo/zero/common/zlog"
	"github.com/SunMaybo/zero/zctl/cmd"
	"github.com/SunMaybo/zero/zctl/download"
	"github.com/SunMaybo/zero/zctl/file"
	"github.com/SunMaybo/zero/zctl/util"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var protocURL = map[string]string{
	"darwin":  "https://github.com/protocolbuffers/protobuf/releases/download/v3.20.0/protoc-3.20.0-osx-x86_64.zip",
	"windows": "https://github.com/protocolbuffers/protobuf/releases/download/v3.20.0/protoc-3.20.0-win64.zip",
	"linux":   "https://github.com/protocolbuffers/protobuf/releases/download/v3.20.0/protoc-3.20.0-linux-x86_64.zip",
	"arm64":   "https://github.com/protocolbuffers/protobuf/releases/download/v3.20.0/protoc-3.20.0-linux-aarch_64.zip",
}
var protoValidateURL = map[string]string{
	"windows": "https://repo1.maven.org/maven2/io/envoyproxy/protoc-gen-validate/protoc-gen-validate/0.6.7/protoc-gen-validate-0.6.7-windows-x86_64.exe",
	"linux":   "https://repo1.maven.org/maven2/io/envoyproxy/protoc-gen-validate/protoc-gen-validate/0.6.7/protoc-gen-validate-0.6.7-linux-x86_64.exe",
	"darwin":  "https://repo1.maven.org/maven2/io/envoyproxy/protoc-gen-validate/protoc-gen-validate/0.6.7/protoc-gen-validate-0.6.7-osx-x86_64.exe",
	"arm64":   "https://repo1.maven.org/maven2/io/envoyproxy/protoc-gen-validate/protoc-gen-validate/0.6.7/protoc-gen-validate-0.6.7-osx-aarch_64.exe",
}

var docURL = map[string]string{
	"windows": "https://github.com/pseudomuto/protoc-gen-doc/releases/download/v1.5.1/protoc-gen-doc_1.5.1_windows_amd64.tar.gz",
	"linux":   "https://github.com/pseudomuto/protoc-gen-doc/releases/download/v1.5.1/protoc-gen-doc_1.5.1_linux_amd64.tar.gz",
	"darwin":  "https://github.com/pseudomuto/protoc-gen-doc/releases/download/v1.5.1/protoc-gen-doc_1.5.1_darwin_amd64.tar.gz",
	"arm64":   "https://github.com/pseudomuto/protoc-gen-doc/releases/download/v1.5.1/protoc-gen-doc_1.5.1_darwin_arm64.tar.gz",
}
var grpcJava = map[string]string{
	"windows": "https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/1.45.1/protoc-gen-grpc-java-1.45.1-windows-x86_64.exe",
	"linux":   "https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/1.45.1/protoc-gen-grpc-java-1.45.1-linux-x86_64.exe",
	"darwin":  "https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/1.45.1/protoc-gen-grpc-java-1.45.1-osx-x86_64.exe",
	"arm64":   "https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/1.45.1/protoc-gen-grpc-java-1.45.1-osx-aarch_64.exe",
}

func GolangInstall() {
	if !cmd.GolangVersionGreaterThan16() {
		zlog.S.Errorw("golang version is less than 1.6")
		return
	}
	if err := cmd.GolangInstallInject(); err != nil {
		zlog.S.Errorw("install protoc-go-inject-tag error", "err", err)
		return
	}
	if err := cmd.GolangInstallGrpc(); err != nil {
		zlog.S.Errorw("install protoc-gen-go-grpc error", "err", err)
		return
	}
	JavaInstall()
}

func JavaInstall() {
	protoPath := DownloadProtoc()
	DownloadProtoValidate()
	DownloadJavaGrpc()
	DownloadProtocDoc()
	_ = filepath.Walk(file.GetFilePath(protoPath, "/bin"), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			err := os.Chmod(path, os.ModePerm)
			if err != nil {
				zlog.S.Error("chmod error:", err)
			}
		}
		return nil
	})

}

func DownloadProtoc() string {
	installKey := getInstallKey()
	if protocURL, ok := protocURL[installKey]; !ok {
		zlog.S.Errorf("%s not support", runtime.GOOS)
		os.Exit(-1)
	} else {
		protoPath := getProtoBasePath()
		dst := file.GetFilePath(getProtoBasePath(), "/protoc.zip")
		if exist, _ := file.PathExists(dst); exist {
			zlog.S.Infof("the protoc already exists. If you need to download it again, delete it yourself from %s", dst)
			return protoPath
		}
		zlog.S.Infof("download protoc from %s", protocURL)
		zlog.S.Infof("if download fail, please download it manually from %s and unzip it to %s", protocURL, protoPath)
		if result, err := downloadFile(protocURL, dst); err != nil {
			zlog.S.Errorf("download protoc error: %s", err)
			os.Exit(1)
		} else {
			zlog.S.Infof("download protoc success: %s", result)
			// unzip
			if err := util.DeCompressZip(result, protoPath); err != nil {
				zlog.S.Errorf("unzip protoc error: %s", err)
				os.Exit(1)
			} else {
				zlog.S.Infof("unzip protoc success")
			}
		}
		return protoPath
	}
	return ""
}

func DownloadProtoValidate() {
	installKey := getInstallKey()
	if protocValidateURL, ok := protoValidateURL[installKey]; !ok {
		zlog.S.Errorf("%s not support", runtime.GOOS)
		os.Exit(-1)
	} else {
		dst := file.GetFilePath(getProtoBasePath(), "/bin/protoc-gen-validate")
		if exist, _ := file.PathExists(dst); exist {
			zlog.S.Infof("the protoc-gen-validate already exists. If you need to download it again, delete it yourself from %s", dst)
			return
		}
		validateProto := file.GetFilePath(getProtoBasePath(), "/include/validate")
		zlog.S.Infof("download protoc-gen-validate from %s", protocValidateURL)
		zlog.S.Infof("if download fail, please download it manually from %s", protocValidateURL)
		if result, err := downloadFile(protocValidateURL, dst); err != nil {
			zlog.S.Errorf("download protoc-gen-validate error: %s", result)
		} else {
			zlog.S.Infof("download protoc-gen-validate success: %s", result)
		}
		file.MkdirAll(validateProto)
		if err := file.WriterFile(file.GetFilePath(validateProto, "/validate.proto"), []byte(validate_067)); err != nil {
			zlog.S.Errorf("save validate.proto err,%s,%s", validateProto, err)
			os.Exit(-1)
		}
	}
}
func DownloadJavaGrpc() {
	installKey := getInstallKey()
	if protocGrpcURL, ok := grpcJava[installKey]; !ok {
		zlog.S.Errorf("%s not support", runtime.GOOS)
		os.Exit(-1)
	} else {
		dst := file.GetFilePath(getProtoBasePath(), "/bin/protoc-gen-java-grpc")
		if exist, _ := file.PathExists(dst); exist {
			zlog.S.Infof("the protoc-gen-java-grpc already exists. If you need to download it again, delete it yourself from %s", dst)
			return
		}
		zlog.S.Infof("download protoc-gen-java-grpc from %s", protocGrpcURL)
		zlog.S.Infof("if download fail, please download it manually from %s", protocGrpcURL)
		if result, err := downloadFile(protocGrpcURL, dst); err != nil {
			zlog.S.Errorf("download protoc-gen-java-grpc error: %s", result)
		} else {
			zlog.S.Infof("download protoc-gen-java-grpc success: %s", result)
		}
	}
}
func DownloadProtocDoc() {
	installKey := getInstallKey()
	basePath := getProtoBasePath()
	if err := file.MkdirAll(basePath); err != nil {
		zlog.S.Errorf("mkdir %s error: %s", basePath, err)
		os.Exit(1)
	}
	if protocURL, ok := docURL[installKey]; !ok {
		zlog.S.Errorf("%s not support", runtime.GOOS)
		os.Exit(-1)
	} else {
		dst := file.GetFilePath(basePath, "/protoc_doc.tar.gz")
		if exist, _ := file.PathExists(dst); exist {
			zlog.S.Infof("the protoc_gen_doc already exists. If you need to download it again, delete it yourself from %s", dst)
			return
		}
		zlog.S.Infof("download protoc_doc from %s", protocURL)
		zlog.S.Infof("if download fail, please download it manually from %s and unzip it to %s", protocURL, basePath)
		if result, err := downloadFile(protocURL, dst); err != nil {
			zlog.S.Errorf("download protoc_doc error: %s", err)
			os.Exit(1)
		} else {
			zlog.S.Infof("download protoc_doc success: %s", result)
			if err := util.UnTarGz(result, file.GetFilePath(basePath, "/bin")+"/"); err != nil {
				zlog.S.Errorf("tarzip protoc_doc error: %s", err)
				os.Exit(1)
			} else {
				zlog.S.Infof("tarzip protoc_doc success:%s", basePath)
			}
		}
	}
}

func downloadFile(url, dst string) (string, error) {
	d := download.NewDownload()
	return d.DownURL(url, dst, func(completeByte, totalSize int64, process float64) {
		zlog.S.Infof("download...... %dMb/%dMb %.2f%%", completeByte/1024/1024, totalSize/1024/1024, process*100)
	})
}
func getProtoBasePath() string {
	protoPath := ""
	homeDir := ""
	if user, err := user.Current(); err != nil {
		zlog.S.Errorf("get current user err:%s", err)
		os.Exit(1)
	} else {
		homeDir = user.HomeDir
	}
	if runtime.GOOS == "windows" {
		protoPath = homeDir + "\\proto"
	} else {
		protoPath = homeDir + "/proto"
	}
	return protoPath
}

func getInstallKey() string {
	installKey := ""
	if runtime.GOARCH == "arm64" && runtime.GOOS == "darwin" {
		installKey = "arm64"
	} else {
		installKey = runtime.GOOS
	}
	return installKey
}
