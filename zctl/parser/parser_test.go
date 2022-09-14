package parser

import "testing"

func TestParser(t *testing.T) {

	_, _ = Parser("../../proto/hello/test_services.proto")
}
func TestParserJavaGrpc(t *testing.T) {
	ParserJavaGrpc("/Users/fico/project/ins-xhwallet-platform/grpc_java/universal/src/main/java/cn/xunhou/grpc/proto/universal/UniversalServiceGrpc.java")
}
func TestParserSQL(t *testing.T) {
	ParserCreatedSQL("/Users/fico/project/java/ins-xhportal-platform/sql")
}
