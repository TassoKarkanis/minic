
minic: build/minic

build/minic: FORCE parser/c_parser.go
	go build -o build/minic minic/minic.go

build/hello: hello/hello.asm
	nasm -f elf64 -o build/hello.o hello/hello.asm
	ld -o build/hello build/hello.o

FORCE:

build/.build:
	mkdir -p build
	touch build/.build

build/.builder: build/.build Dockerfile-builder
	docker build -t minic-builder:latest -f Dockerfile-builder .
	touch build/.builder

parser/c_parser.go: C.g4 build/.builder
	mkdir -p parser
	chmod g+s parser
	docker run \
	    	-v `pwd`:/w \
		-w /w \
		minic-builder:latest \
		antlr4 -Dlanguage=Go -o parser C.g4

tests: build/.test-empty-function

build/.test-empty-function: build/minic tests/empty-function.c
	build/minic tests/empty-function.c
	touch build/.test-empty-function
