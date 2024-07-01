generate:
	rm -rf internal/pkg/pb
	mkdir -p internal/pkg/pb

	protoc \
		--proto_path=api/ \
		--go_out=internal/pkg/pb \
		--go-grpc_out=internal/pkg/pb \
		api/*.proto
