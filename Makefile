
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

TESTS_C := $(shell find tests -name "*.c")
TESTS_ASM := $(TESTS_C:tests/%.c=build/%.asm)
TESTS_DIFF := $(TESTS_C:tests/%.c=build/%.diff)

# $(info $$TEST_TARGETS is [${TEST_TARGETS}])

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

build/%.asm : tests/%.c build/minic
	build/minic -S $< -o $@ > $(@:.asm=.out)

.PRECIOUS: build/%.asm

build/%.gcc.asm : tests/%.c
	gcc -fno-asynchronous-unwind-tables -S -O2 $< -o $@

build/%.diff : tests/%.asm build/%.asm
	diff -c $< $(@:.diff=.asm) > $@

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

# objconv -v0 -fnasm $(@:.asm=-gcc.o) $(@:.asm=-gcc.asm)
