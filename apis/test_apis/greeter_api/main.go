package main

import (
	"flag"
	"github.com/LasercatDAO/cat-shop/svc"
	"github.com/LasercatDAO/cat-shop/zlog"
	"log"
	"os"

	"github.com/LasercatDAO/cat-shop/etc"
	"github.com/LasercatDAO/cat-shop/server"
)

// 设置路由信息
var path string

func init() {
	flag.StringVar(&path, "etc", "./etc/config.yaml", "config for filepath")
	flag.Parse()
}
func main() {
	if cfg, err := etc.New(path); err != nil {
		log.Fatal(err)
	} else {
		if cfg.FileDir == "" {
			cfg.FileDir = "./assets"
		}
		os.MkdirAll(cfg.FileDir, 0777)
		zlog.InitLogger(false)
		svc := svc.NewServiceContext(cfg)
		svc.Cfg.FileDir = svc.Cfg.FileDir + "/images"
		os.MkdirAll(svc.Cfg.FileDir, 0777)
		server.NewServer(svc).Start()
	}
}
