
all: generate build

generate:
	${GOPATH}/bin/gogdc generate --path="./src"

build:
	go build -v -buildmode=c-shared -o ./bin/libsimple.so ./src/*.go

clean:
	go clean
	rm ./bin/libsimple.so
	rm ./src/*_registrable.gen.go

.PHONY: all
