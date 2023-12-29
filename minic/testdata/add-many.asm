section .data

section .text
	global f

f:
	; param 'int' a0 -> edi [rsp - 0]
	; param 'int' a1 -> esi [rsp - 4]
	; param 'int' a2 -> edx [rsp - 8]
	; param 'int' a3 -> ecx [rsp - 12]
	; param 'int' a4 -> r8d [rsp - 16]
	; param 'int' a5 -> r9d [rsp - 20]
	push rbx
	push r12
	push r13
	push r14
	push r15
	mov eax, edi
	add eax, esi
	mov ebx, eax
	add ebx, edx
	mov eax, ebx
	add eax, ecx
	mov ebx, eax
	add ebx, r8d
	mov eax, ebx
	add eax, r9d
	jmp f.end
f.end:
	pop r15
	pop r14
	pop r13
	pop r12
	pop rbx
	ret

