proto_gen:
	protoc \
		--go_out=gen --go_opt=paths=source_relative \
		--go-grpc_out=gen --go-grpc_opt=paths=source_relative \
		--proto_path=idl/proto "idl/proto/helloworld.proto"