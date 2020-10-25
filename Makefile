
minic: build/minic

build/minic: FORCE parser/c_parser.go
	go build -o build/minic ./minic

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

tests: build/empty-function.asm

build/empty-function.asm: build/minic tests/empty-function.c
	build/minic -S tests/empty-function.c -o build/empty-function.asm
	nasm -f elf64 -o build/empty-function.o build/empty-function.asm
	gcc -fno-asynchronous-unwind-tables -c -O2 tests/empty-function.c -o build/empty-function-gcc.o
	objconv -fnasm build/empty-function-gcc.o > build/empty-function-gcc.asm
