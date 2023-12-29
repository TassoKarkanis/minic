minic: build/minic

SOURCES := go.mod go.sum \
    $(shell find codegen -name "*.go") \
    $(shell find minic -name "*.go") \
    $(shell find parser -name "*.go") \
    $(shell find symbols -name "*.go") \
    $(shell find types -name "*.go")
build/minic: $(SOURCES) parser/c_parser.go build/.build
	go build -o build/minic ./minic

build/hello: hello/hello.asm
	nasm -f elf64 -o build/hello.o hello/hello.asm
	ld -o build/hello build/hello.o

parser/c_parser.go: C.g4 build/.builder
	mkdir -p parser
	chmod g+s parser
	docker run \
	    	-v `pwd`:/w \
		-w /w \
		minic-builder:latest \
		antlr4 -Dlanguage=Go -o parser C.g4

clean:
	rm -f minic/testdata/*.out.asm || true
	rm -f minic/testdata/*.o || true
	rm -f minic/testdata/*.out || true

devcontainer: build/.devcontainer FORCE
	docker kill minic-devcontainer || true
	docker rm minic-devcontainer || true
	docker run \
		--name minic-devcontainer \
		--network agmt \
		-d \
		-e CGO_ENABLED=0 \
		-v `pwd`:/w \
		-v minic-buildvol:/root/go \
		-w /w \
		minic-devcontainer:latest \
		sleep infinity

test: build/.builder
	docker run \
		-v minic-buildvol:/root/go \
		-v `pwd`:/w \
		-w /w \
		-it minic-builder:latest \
		go test -v ./minic/...

shell: build/.builder
	docker run \
		-v minic-buildvol:/root/go \
		-v `pwd`:/w \
		-w /w \
		-it minic-builder:latest \
		bash

tests: $(TESTS_DIFF)

build/.devcontainer: build/.build Dockerfile
	docker build \
		-t minic-devcontainer:latest \
		--target minic-devcontainer \
		-f Dockerfile \
		.
	touch build/.devcontainer

build/.builder: build/.build Dockerfile
	docker build \
		-t minic-builder:latest \
		--target minic-builder \
		-f Dockerfile \
		.
	touch build/.builder

build/.build:
	mkdir -p build
	touch build/.build

FORCE: