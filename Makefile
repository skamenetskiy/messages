.PHONY: build
build:
	go build -o bin/messages internal/cmd/main.go

.PHONY: build-static
build-static:
	CGO_ENABLED=0 go build -a -tags netgo -ldflags="-w -s -extldflags '-static'" -o bin/messages internal/cmd/main.go

.PHONY: proto
proto:
	protoc \
		--proto_path api \
		--go_out=api \
		--go_opt=paths=source_relative \
    	--go-grpc_out=api \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out api \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		--grpc-gateway_opt allow_delete_body=true \
		--openapiv2_out api \
		--openapiv2_opt allow_delete_body=true \
    	api/messages.proto