package main

import (
	"flag"
	"github.com/SunMaybo/zero/zctl/gen"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var serviceType = flag.String("t", "services", "service type services or apis")

//var serviceProtocol = flag.String("s", "rpc", "service protocol rpc or http")
var project = flag.String("p", "", "project name")

func main() {
	flag.Parse()
	if *project == "" {
		log.Fatal("project name is empty")
	}
	workDir := getCurrentAbsolutePath()
	b := gen.NewRpcBuilder(getProjectByMod(), workDir, *project, *serviceType)
	b.StartBuild()
}
func getCurrentAbsolutePath() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}
func getProjectByMod() string {
	buff, err := ioutil.ReadFile("go.mod")
	if err != nil {
		log.Printf("read go.mod error:%v", err)
		log.Fatalln(err)
	}
	str := string(buff)
	for _, s := range strings.Split(str, "\n") {
		if strings.HasPrefix(strings.TrimSpace(s), "module") {
			return strings.TrimSpace(strings.Split(s, " ")[1])
		}
	}
	return ""
}
