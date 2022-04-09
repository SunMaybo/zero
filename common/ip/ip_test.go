package ip

import (
	"fmt"
	"testing"
)

func TestLower16BitPrivateIP(t *testing.T) {
	ip, _ := PrivateIPv4()
	fmt.Println(ip.String())
}
func TestLocalHostIP(t *testing.T) {
	fmt.Println(LocalHostIP())
}
