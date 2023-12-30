section .data

section .text
	global f

f:
	; param 'int' a -> edi [rsp - 0]
	; param 'int' b -> esi [rsp - 4]
	; param 'int' c -> edx [rsp - 8]
	push rbx
	push r12
	push r13
	push r14
	push r15
	mov dword [rsp - 8], edx ; flush to local
	mov eax, esi
	mul edi
	; var 'int' x -> ebx [rsp - 16]
	mov ebx, eax ; update l-value register
	mov dword [rsp - 16], ebx ; store l-value
	mov eax, ebx
	add eax, dword [rsp - 8]
	jmp f.end
f.end:
	pop r15
	pop r14
	pop r13
	pop r12
	pop rbx
	ret

