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
	mov eax, edi ; greater-equal: load LHS
	cmp eax, esi ; greater-equal: compare
	setge al ; greater-equal: set byte in result
	movzx eax, al ; greater-equal: zero-extend
	jmp f.end
f.end:
	pop r15
	pop r14
	pop r13
	pop r12
	pop rbx
	ret

