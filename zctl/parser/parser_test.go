package parser

import "testing"

func TestParser(t *testing.T) {

	_, _ = Parser("../../proto/hello/test_services.proto")
}
func TestParserJavaGrpc(t *testing.T) {
	ParserJavaGrpc("/Users/fico/go/src/zero/grpc_java/test_apis/src/main/java/com/xuhou/middle/proto/test_apis/GreeterGrpc.java")
}
