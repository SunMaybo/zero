package gen

import (
	"github.com/SunMaybo/zero/common/zlog"
	"github.com/SunMaybo/zero/zctl/cmd"
)

func GenDoc(src, dst string, docType int) {
	if err := cmd.GetProtoDoc(src, dst, cmd.DocType(docType)); err != nil {
		zlog.S.Errorw("gen doc error", "err", err)
		return
	}
}
