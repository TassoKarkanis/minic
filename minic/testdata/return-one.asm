section .data

section .text
	global f

f:
	mov eax, 1
	jmp f.end
f.end:
	ret

