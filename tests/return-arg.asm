section .data

section .text
	global f

f:
	mov eax, 1
	jmp blah.end
blah.end:
	ret

