services=''
rpc:
	protoc --proto_path=proto --go_out=plugins=grpc:proto/test  proto/$(services)/greeter.proto
rpc_java:
	protoc -I=/usr/local/include/google/protobuf -I=. --java_out=. greeter.proto
	protoc --plugin=protoc-gen-grpc-java=/usr/local/bin/protoc-gen-java --plugin=protoc-gen-validate=/usr/local/bin/protoc-gen-validate -I=/usr/local/include/google/protobuf -I=. --validate-java_out=. --java_out=. --grpc-java_out=. *.proto
	protoc --plugin=protoc-gen-grpc-java=/usr/local/bin/protoc-gen-java --plugin=protoc-gen-validate-java=/usr/local/bin/protoc-gen-validate -I=/usr/local/include/google/protobuf -I=. --validate_out="lang=java:." --java_out=. --grpc-java_out=. *.proto
