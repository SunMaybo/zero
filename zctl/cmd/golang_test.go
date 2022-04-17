package cmd

import (
	"github.com/SunMaybo/zero/common/zlog"
	"testing"
)

func init() {
	zlog.InitLogger(false)
}

func TestGoVersion(t *testing.T) {
	t.Log(GolangVersionGreaterThan16())
}
func TestGetGoPath(t *testing.T) {
	t.Log(GetGolangBinPath())
}
func TestGolangInstallInject(t *testing.T) {
	t.Log(GolangInstallInject())
}
func TestGolangInstallGrpc(t *testing.T) {
	t.Log(GolangInstallGrpc())
}
