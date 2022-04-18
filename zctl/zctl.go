package main

import (
	"github.com/SunMaybo/zero/common/zcfg"
	"github.com/SunMaybo/zero/common/zlog"
	"github.com/SunMaybo/zero/zctl/config"
	"github.com/SunMaybo/zero/zctl/execute"
	"github.com/SunMaybo/zero/zctl/file"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"os/user"
)

func main() {
	u, _ := user.Current()
	cfgPath := file.GetFilePath(u.HomeDir, "/.zctl")
	if exist, err := file.PathExists(cfgPath); err != nil || !exist {
		if err := file.MkdirAll(cfgPath); err != nil {
			panic(err)
		}
	}
	cfgFilePath := file.GetFilePath(u.HomeDir, "/.zctl/config.yaml")
	if exist, err := file.PathExists(cfgFilePath); err != nil || !exist {
		genConfig(cfgFilePath)
	}
	cfg := config.Config{}
	zcfg.LoadConfig(cfgFilePath, &cfg)
	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(execute.GetAllCommands(cfg)...)
	rootCmd.Execute()

}
func genConfig(path string) {
	cfg := config.Config{
		Maven:                     "/usr/local/maven/bin/mvn",
		MavenDeploymentRepository: "releases-snapshots::default::zapi://nexus.tongdao.cn/nexus/content/repositories/releases-snapshots/",
		MavenSettings:             file.GetFilePath(getHomeDir(), "/.m2/settings.xml"),
	}
	buff, _ := yaml.Marshal(&cfg)
	if err := file.WriterFile(path, buff); err != nil {
		zlog.S.Errorf("write config file error:%s", err.Error())
		os.Exit(0)
	}
}
func getHomeDir() string {
	if u, err := user.Current(); err != nil {
		panic(err)
	} else {
		return u.HomeDir
	}
}
