
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

C_TESTS := $(shell find tests -name "*.c")
$(info $$C_TESTS is [${C_TESTS}])

TEST_TARGETS := $(C_TESTS:tests/%.c=build/%.asm)
$(info $$TEST_TARGETS is [${TEST_TARGETS}])

tests: $(TEST_TARGETS)

build/%.asm : tests/%.c build/minic
	build/minic -S $< -o $@
	@diff -c $(<:.c=.asm) $@
	@nasm -f elf64 -o $(@:.asm=.o) $@
	@gcc -fno-asynchronous-unwind-tables -c -O2 $< -o $(@:.asm=-gcc.o)
	@objconv -v0 -fnasm $(@:.asm=-gcc.o) $(@:.asm=-gcc.asm)
