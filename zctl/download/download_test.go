package download

import (
	"fmt"
	"testing"
)

func TestDownTask_DownURL(t *testing.T) {
	downloadTask := NewDownload()
	err, resp := downloadTask.DownURL("https://epidentity.oss-cn-shanghai.aliyuncs.com/pipeline-assistant/Docker_Installer.exe", "/Users/mayunbao/Desktop/", func(completeByte, totalSize int64, process float64) {
		fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
			completeByte,
			totalSize,
			100*process)

	})
	fmt.Println("out:", err)
	fmt.Println("xxx", resp)

}
