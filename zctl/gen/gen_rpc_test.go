package gen

import "testing"

func TestGenRpc(t *testing.T) {
	b := NewRpcBuilder("zero", "/Users/fico/go/src/zero", "test_services", "services")
	b.StartBuild()
}
