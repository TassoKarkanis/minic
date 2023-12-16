section .data

section .text
	global f

f:
	mov eax, 1
	add eax, 1
	mov dword [rsp - 4], eax
	mov eax, dword [rsp - 4]
	jmp blah.end
blah.end:
	ret

