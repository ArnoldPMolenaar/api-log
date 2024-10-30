APP_NAME = log-api
BUILD_DIR = $(PWD)/build

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

lint:
	golangci-lint run ./...
