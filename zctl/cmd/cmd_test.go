package cmd

import "testing"

func TestMavenVersion(t *testing.T) {
	t.Log(MavenExist("/usr/local/maven/bin/mvn"))

}
