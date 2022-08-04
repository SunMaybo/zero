package gen

import (
	"testing"
)

func TestGenProto(t *testing.T) {
	if err := GenerateSchema("/Users/fico/Desktop/project_pc/ins-xhportal-platform/sql", "test"); err != nil {
		t.Fatal(err)
	}

}
func TestGenEntity(t *testing.T) {
	if err := GenerateDao("/Users/fico/Desktop/project_pc/ins-xhportal-platform/sql", "test"); err != nil {
		t.Fatal(err)
	}
}
