package install

import (
	"testing"
)

func TestInstall(t *testing.T) {
	DownloadProtocDoc()
}
func TestDownloadJavaGrpc(t *testing.T) {
	DownloadJavaGrpc()
}
func TestJavaInstall(t *testing.T) {
	JavaInstall()
}
func TestGolangInstall(t *testing.T) {
	GolangInstall()
}
