.PHONY: proto test
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/processor.proto

mocks:
	mockery

test:
	go test -v ./... -covermode=atomic -coverprofile=coverage.out

dapr_run:
	dapr run --app-id=processing-orchestrator --app-port 55555 --dapr-grpc-port 50011 --resources-path ./dapr/components -- go run cmd/server.go

dapr:
	dapr run --app-id=processing-orchestrator --app-protocol=grpc --app-port 55555 --dapr-grpc-port 50011  --resources-path ./dapr/components

container:
	docker build -t processing-orchestrator .


# Waiting for dapr multi-app run...
run_e2e:
	docker run -p -d 9000:9000 -p 9001:9001 minio/minio server /data --console-address ":9001"
	# Cooking server
	cd ../Pandora-cooking-server/ && dapr run --log-level debug --dapr-http-max-request-size 300 --app-id cooking-server --app-port 3004 --dapr-http-port 3500 --resources-path ./dapr/components -- npm run start:dev && cd -
	# Encode box
	cd ../encode-box/ && dapr run --log-level debug --app-id encode-box --dapr-http-max-request-size="1000" --dapr-http-port=3500 --dapr-grpc-port=50010 --components-path=dapr/components && cd -