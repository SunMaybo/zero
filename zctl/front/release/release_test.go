package release

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"testing"
)

func TestDecrypt(t *testing.T) {
	Delay("format", ".", "", false, "", "", "", true, false)
}
func TestEncrypt(t *testing.T) {
	t.Log(EncryptByAes("xbb123:::", []byte("")))
}

func TestUploadCompanyWebsite(t *testing.T) {
	endpoint := "oss-cn-beijing.aliyuncs.com"
	client, err := oss.New(endpoint, "", "")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 填写存储空间名称，例如examplebucket。
	bucket, err := client.Bucket("xbanban-site")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	uploadDirectoryFileTree(bucket, "/Users/fico/Desktop/dist/", "")
}
