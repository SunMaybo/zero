package gen

import (
	"io/fs"
	"path/filepath"
	"testing"
)

func TestGenJAVARpc(t *testing.T) {
	err := filepath.Walk("/Users/fico/go/src/zero/grpc_java/test_apis/src/main/java/com/xuhou/middle/proto/test_apis", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		return javaGrpcImpl("/Users/fico/go/src/zero/grpc_java/test_apis/src/main/java/com/xuhou/middle/proto/test_apis", path)
	})
	if err != nil {
		t.Fatal(err)
	}
}
