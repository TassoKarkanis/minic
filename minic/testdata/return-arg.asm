section .data

section .text
	global f

f:
	; param 'int' x -> edi
	mov eax, edi
	jmp f.end
f.end:
	ret

