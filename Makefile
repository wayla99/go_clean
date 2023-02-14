run:
	go run main.go

swagger:
	swag init --dir ./src/interface/fiber_server --output ./src/interface/fiber_server/docs

unit-test:
	go test -v ./src/...
