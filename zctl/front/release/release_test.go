package release

import (
	"testing"
)

func TestDecrypt(t *testing.T) {
	Delay("format", ".", "", false, "", "", "", true, false)
}
func TestEncrypt(t *testing.T) {
	t.Log(EncryptByAes("xbb123:::", []byte("LTAI5t89dnZgsEC99idBKb67-5eYCj3sKYLeyScX0kcfUrKDJ1kUIpB")))
}
