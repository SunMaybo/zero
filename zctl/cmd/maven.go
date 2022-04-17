package cmd

import (
	"errors"
	"fmt"
	"github.com/SunMaybo/zero/common/zlog"
	"log"
	"strings"
)

func MavenDeploy(maven, mavenSettings, altDeploymentRepository, dir string) error {
	zlog.S.Infof("Deploying maven project %s,cmd:%s", dir, maven+" clean deploy -DaltDeploymentRepository="+altDeploymentRepository)
	if !MavenExist(maven) {
		return errors.New("maven not exist")
	}
	isOk := false
	result := ""
	if err := Execute(maven, dir, func(lines string) {
		if strings.Contains(lines, "BUILD SUCCESS") {
			isOk = true
		}
		log.Println(lines)
	}, "clean", "deploy", "-gs="+mavenSettings, fmt.Sprintf("-DaltDeploymentRepository=%s", altDeploymentRepository)); err != nil {
		return err
	}
	log.Println(result)
	if !isOk {
		return errors.New("maven deploy failed")
	}
	return nil
}
func MavenExist(maven string) bool {
	isOk := false
	if err := Execute(maven, "", func(lines string) {
		if strings.Contains(strings.ToLower(lines), "maven") {
			isOk = true
		}
	}, "-v"); err != nil {
		log.Println(err)
		return false
	}

	return isOk
}
