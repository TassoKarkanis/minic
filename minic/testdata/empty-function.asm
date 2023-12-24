section .data

section .text
	global f

	mov dword [rsp - 4], edi
f.end:
	ret

