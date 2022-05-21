gen:
	protoc --go_out=./ --go-grpc_out=./ --proto_path=proto proto/*.proto

clean:
	rm ./protobuf/*

run:
	go run main.go