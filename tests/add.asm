section .data

section .text
	global f

f:
	mov r9d, 1
	add r9d, 1
	mov eax, r9d
	ret
