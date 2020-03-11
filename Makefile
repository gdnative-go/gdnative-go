
all: clean generate build

clean:
	rm gdnative/*.gen.* || true
	go clean

generate:
	go generate

build:
	go build -x ./gdnative/...

.PHONY: all
