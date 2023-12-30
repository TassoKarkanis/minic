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
	mov eax, edi
	add eax, 1
	; var 'int' x -> ebx [rsp - 12]
	mov ebx, eax ; update l-value register
	mov dword [rsp - 12], ebx ; store l-value
	; var 'int' y -> eax [rsp - 16]
	mov r12d, esi
	add r12d, 1
	; var 'int' x -> r13d [rsp - 24]
	mov r13d, r12d ; update l-value register
	mov dword [rsp - 24], r13d ; store l-value
	mov eax, r13d ; update l-value register
	mov dword [rsp - 16], eax ; store l-value
	mov r12d, ebx
	add r12d, eax
	mov eax, r12d
	jmp f.end
f.end:
	pop r15
	pop r14
	pop r13
	pop r12
	pop rbx
	ret

