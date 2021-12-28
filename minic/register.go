package main

// Registers and Calling Convention
//
// 1) Caller-save (value may be different on return from call):
//    rax, rcx, rdx, rdi, rsi, rsp, r8-r11
//
// 2) Callee-save (value preserved on return from call):
//    rbx, rbp, r12-15
//
// 3) First six integer/address parameters are stored in:
//    rdi, rsi, rdx, rcx, r8, r9
//
// 4) Stack grows down.
//
//

type Register struct {
	integer   bool // as opposed to floating point
	fullName  string
	dwordName string
	wordName  string
	byteName  string
	binding   *Value
}

func (r Register) Name(byteSize int) string {
	switch byteSize {
	case 8:
		return r.fullName

	case 4:
		return r.dwordName

	case 2:
		return r.wordName

	case 1:
		return r.byteName

	default:
		panic("invalid bytesize!")
	}
}
