section .data
    msg db      "hello, world!", 10

section .text
	global _start
	global _write
	global _exit
	global _strlen
_start:
	;; print
	mov rax, msg
	call _print_hello
	call _exit

;; int write(int fd, const void *buf, size_t count)
_print_hello:
	mov     rax, 1
	mov	rdi, 1
	mov     rsi, msg
	mov     rdx, 14	
	syscall
	ret

;; int write(int fd, const void *buf, size_t count)
_write:
	

;; void exit(int rc)
_exit:
	mov	rdi, rax
	mov	rax, 60	
	mov	rdi, 0
	syscall

;; void strlen(const char* s)
_strlen:
	;; Use rbx to step through string.  Keep rax to perform subtraction
	;; at the end.
	mov rbx, rax
	dec rbx

_strlen_increment_and_check_zero:
	inc rbx
	movsx rcx, byte [rbx]
	cmp rcx, 0
	jz _strlen_end
	jmp _strlen_increment_and_check_zero

_strlen_end:
	mov rax, rbx
	ret
	
