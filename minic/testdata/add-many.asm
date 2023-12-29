section .data

section .text
	global f

f:
	; param 'int' a0 -> edi
	; param 'int' a1 -> esi
	; param 'int' a2 -> edx
	; param 'int' a3 -> ecx
	; param 'int' a4 -> r8d
	; param 'int' a5 -> r9d
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

