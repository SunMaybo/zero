package release

import (
	"fmt"
	"testing"
)

func TestRelease(t *testing.T) {
	Delay("qa39", "/Users/fico/project/java/ins-xhwallet-platform", true)
}
func TestDecrypt(t *testing.T) {
	pwd, _ := EncryptByAes("xbb123:::", []byte("LTAI5tQAMk47JoGssQmSq9E6-69IYL0TcAMdX3HLVLvQPWApqbkhHry-xbbossuploader"))
	fmt.Println(string(pwd))
	buff, _ := DecryptByAes("xbb123:::", pwd)
	fmt.Println(string(buff))
}
