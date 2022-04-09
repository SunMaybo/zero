services=''
rpc:
	protoc --proto_path=proto --go_out=plugins=grpc:proto/test  proto/$(services)/greeter.proto
