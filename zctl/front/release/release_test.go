package release

import (
	"testing"
)

func TestDecrypt(t *testing.T) {
	Delay("format", ".", "", false, "", "", "", true, false)
}
func TestEncrypt(t *testing.T) {
	t.Log(EncryptByAes("xbb123:::", []byte("")))
}
