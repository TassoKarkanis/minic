section .data

section .text
	global f

f:
	mov edi, 1
	add edi, 1
	mov eax, edi
	jmp f.end
f.end:
	ret

