section .data

section .text
	global f

f:
	; param 'int' a -> edi [rsp - 0]
	; param 'int' b -> esi [rsp - 4]
	push rbx
	push r12
	push r13
	push r14
	push r15
	mov rdx, 0
	mov eax, edi
	cdq
	idiv esi
	mov eax, edx
	jmp f.end
f.end:
	pop r15
	pop r14
	pop r13
	pop r12
	pop rbx
	ret

