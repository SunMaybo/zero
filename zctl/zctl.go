package main

import (
	"flag"
	"fmt"
	"github.com/SunMaybo/zero/zctl/gen"
	"os"
	"strings"
)

var serviceType = flag.String("t", "services", "service type services or apis")
var serviceProtocol = flag.String("s", "rpc", "service protocol rpc or http")
var project = flag.String("p", "", "project name")

func main() {
	flag.Parse()
	if *project == "" {
		fmt.Println("project name is empty")
		os.Exit(1)
	}
	workDir := getCurrentAbsolutePath()
	p := strings.Split(workDir, "/")[len(strings.Split(workDir, "/"))-1]
	b := gen.NewRpcBuilder(p, workDir, *project, *serviceType)
	b.StartBuild()
}
func getCurrentAbsolutePath() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}
