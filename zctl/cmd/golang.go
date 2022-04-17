package cmd

import (
	"errors"
	"github.com/SunMaybo/zero/common/zlog"
	"os"
	"strconv"
	"strings"
)

func GolangVersionGreaterThan16() bool {
	if result, err := Run("go version", ""); err != nil {
		zlog.S.Errorf("go version error:%s", err)
		return false
	} else if !strings.Contains(result, "go version") {
		return false
	} else {
		version := result[strings.Index(result, "go1.")+4 : strings.Index(result, "go1.")+6]
		v, err := strconv.ParseInt(version, 10, 64)
		if err != nil {
			zlog.S.Errorf("go version error:%s", err)
			return false
		}
		if v > 16 {
			zlog.S.Info(result)
			return true
		} else {
			zlog.S.Errorf("go version error:%s", "go version must >=1.7")
			return false
		}
	}
}

func GetGolangBinPath() (string, error) {
	path := os.Getenv("GOPATH")
	if path == "" {
		return "", errors.New("GOPATH is empty")
	}
	if strings.HasSuffix(path, "/src") {
		path = path[:len(path)-4]
	}
	return path + "/bin", nil
}

//go install github.com/favadi/protoc-go-inject-tag@latest
//go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

func GolangInstallInject() error {
	result, err := Run("go install github.com/favadi/protoc-go-inject-tag@latest", "")
	if err != nil {
		return err
	} else {
		zlog.S.Info(result)
		return nil
	}
}
func GolangInstallGrpc() error {
	result, err := Run("go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest", "")
	if err != nil {
		return err
	} else {
		zlog.S.Info(result)
		return nil
	}
}
