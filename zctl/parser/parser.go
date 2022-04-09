package parser

import (
	"errors"
	"io/ioutil"
	"strings"
)

type RpcMetadata struct {
	PackageName string
	ServiceName string
	MethodSigns []MethodSign
}
type MethodSign struct {
	Name                string
	Param               string
	IsStreamParam       bool
	ReturnParam         string
	IsStreamReturnParam bool
}

func Parser(filePath string) (RpcMetadata, error) {
	md := RpcMetadata{}
	if buff, err := ioutil.ReadFile(filePath); err != nil {
		return md, err
	} else {
		content := string(buff)
		items := strings.Split(content, "\n")
		for _, item := range items {

			//解析服务名称
			if strings.Contains(item, "go_package") {
				end := strings.LastIndex(item, "\"")
				start := strings.Index(item, "\"")
				if start >= end || start < 0 {
					panic(errors.New("proto is not exist go_package,file:" + filePath))
				}
				md.PackageName = strings.ReplaceAll(item[start+1:end], "proto/", "")
				md.PackageName = strings.ReplaceAll(md.PackageName, "/", "")
			}
			//解析服务名称
			if strings.HasPrefix(strings.TrimSpace(item), "service") {
				start := strings.Index(item, "service")
				end := strings.Index(item[start:], "{")
				if start >= end || start < 0 {
					return md, errors.New("proto is not exist go_package,file:" + filePath)
				}
				md.ServiceName = strings.TrimSpace(item[start+7 : start+end])
			}
			//解析方法签名
			if strings.HasPrefix(strings.TrimSpace(item), "rpc") {
				ms := MethodSign{}
				//解析方法名称
				start := strings.Index(item, "rpc")
				end := strings.Index(item[start:], "(")
				ms.Name = strings.TrimSpace(item[start+3 : start+end])
				//解析方法参数
				itt := item[start+end:]
				idxStart := strings.Index(itt, "(")
				idxEnd := strings.Index(itt[idxStart:], ")")
				param := strings.TrimSpace(itt[idxStart+1 : idxStart+idxEnd])
				if strings.Contains(param, "stream") {
					ms.IsStreamParam = true
					param = strings.TrimSpace(strings.ReplaceAll(param, "stream", ""))
				}
				ms.Param = param

				//解析返回参数
				rit := itt[idxStart+idxEnd:]
				ritStart := strings.Index(rit, "(")
				ritEnd := strings.Index(rit[ritStart:], ")")
				rParam := strings.TrimSpace(rit[ritStart+1 : ritStart+ritEnd])
				if strings.Contains(rParam, "stream") {
					ms.IsStreamReturnParam = true
					rParam = strings.TrimSpace(strings.ReplaceAll(rParam, "stream", ""))
				}
				ms.ReturnParam = rParam

				md.MethodSigns = append(md.MethodSigns, ms)

			}

		}
	}
	return md, nil
}
