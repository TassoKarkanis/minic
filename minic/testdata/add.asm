section .data

section .text
	global f

	mov dword [rsp - 4], edi
	mov eax, 1
	add eax, 1
	mov dword [rsp - 8], eax
	mov eax, dword [rsp - 8]
	jmp f.end
f.end:
	ret

