all: update_buf_deps generate_buf update_go_deps format

prepare:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install go.einride.tech/aip/cmd/protoc-gen-go-aip@v0.66.0

update_buf_deps:
	buf mod update

update_go_deps:
	go mod tidy

generate_buf:
	buf generate

format:
	go fmt ./...
	gofumpt -l -w .

clean_proto_go:
	rm api/proto/*.go
