package download

import (
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"os"
	"time"
)

type DownTask struct {
}

func NewDownload() (dt *DownTask) {
	return &DownTask{}
}

// DownURL .
func (d *DownTask) DownURL(URL, dst string, onProcess func(completeByte, totalSize int64, process float64)) (string, error) {
	// create client
	client := grab.NewClient()
	req, err := grab.NewRequest(dst, URL)
	if err != nil {
		return "", err
	}
	// start download

	resp := client.Do(req)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			onProcess(resp.BytesComplete(), resp.Size(), resp.Progress())
		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		return "", err
	}
	return resp.Filename, nil
}
