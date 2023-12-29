section .data

section .text
	global f

f:
	; param 'int' a0 -> edi [rsp - 0]
	; param 'int' a1 -> esi [rsp - 4]
	push rbx
	push r12
	push r13
	push r14
	push r15
	mov eax, edi
	sub eax, esi
	jmp f.end
f.end:
	pop r15
	pop r14
	pop r13
	pop r12
	pop rbx
	ret

