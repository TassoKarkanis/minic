section .data

section .text
	global f

f:
	mov eax, edi
	jmp f.end
f.end:
	ret

