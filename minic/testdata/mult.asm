section .data

section .text
	global f

f:
	; param 'int' a0 -> edi
	; param 'int' a1 -> esi
	push rbx
	push r12
	push r13
	push r14
	push r15
	mov eax, esi
	mul edi
	jmp f.end
f.end:
	pop r15
	pop r14
	pop r13
	pop r12
	pop rbx
	ret

