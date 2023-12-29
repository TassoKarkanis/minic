package codegen

// Registers and Calling Convention
//
// 1) Caller-save (value may be different on return from call):
//    rax, rcx, rdx, rdi, rsi, rsp, r8-r11
//
// 2) Callee-save (value must be preserved on return from call):
//    rbx, rbp, r12-15
//
// 3) First six integer/address parameters are stored in:
//    rdi, rsi, rdx, rcx, r8, r9
//
// 4) Stack grows down.
//

const RAX = "rax"
const RBX = "rbx"
const RDI = "rdi"
const RSI = "rsi"
const RDX = "rdx"
const RCX = "rcx"
const R8 = "r8"
const R9 = "r9"
const R10 = "r10"
const R11 = "r11"
const R12 = "r12"
const R13 = "r13"
const R14 = "r14"
const R15 = "r15"

// const
var integerParameterOrder = [...]string{
	RDI, RSI, RDX, RCX, R8, R9,
}

// const
var integerRegisterAllocationOrder = [...]string{
	RAX,                     // not needed unil return, and may well contain result
	RBX, R12, R13, R14, R15, // we always save these anyway
	R9, R8, RCX, RDX, RSI, RDI, // in reverse parameter order
}

type Register struct {
	integer    bool // as opposed to floating point
	callerSave bool
	fullName   string
	dwordName  string
	wordName   string
	byteName   string
	binding    *Value
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

func AllocateIntegerRegisters() map[string]*Register {
	// make the map of general-purpose registers
	registers := map[string]*Register{
		RAX: {
			integer:    true,
			callerSave: true,
			fullName:   "rax",
			dwordName:  "eax",
			wordName:   "ax",
			byteName:   "al",
		},
		RBX: {
			integer:    true,
			callerSave: false,
			fullName:   "rbx",
			dwordName:  "ebx",
			wordName:   "bx",
			byteName:   "bl",
		},
		RDI: {
			integer:    true,
			callerSave: true,
			fullName:   "rdi",
			dwordName:  "edi",
			wordName:   "di",
			byteName:   "dl",
		},
		RSI: {
			integer:    true,
			callerSave: true,
			fullName:   "rsi",
			dwordName:  "esi",
			wordName:   "si",
			byteName:   "sl",
		},
		RDX: {
			integer:    true,
			callerSave: true,
			fullName:   "rdx",
			dwordName:  "edx",
			wordName:   "dx",
			byteName:   "dl",
		},
		RCX: {
			integer:    true,
			callerSave: true,
			fullName:   "rcx",
			dwordName:  "ecx",
			wordName:   "cx",
			byteName:   "cl",
		},
		R8: {
			integer:    true,
			callerSave: true,
			fullName:   "r8",
			dwordName:  "r8d",
			wordName:   "r8w",
			byteName:   "r8b",
		},
		R9: {
			integer:    true,
			callerSave: true,
			fullName:   "r9",
			dwordName:  "r9d",
			wordName:   "r9w",
			byteName:   "r9b",
		},
		R10: {
			integer:    true,
			callerSave: true,
			fullName:   "r10",
			dwordName:  "r10d",
			wordName:   "r10w",
			byteName:   "r10b",
		},
		R11: {
			integer:    true,
			callerSave: true,
			fullName:   "r11",
			dwordName:  "r11d",
			wordName:   "r11w",
			byteName:   "r11b",
		},
		R12: {
			integer:    true,
			callerSave: false,
			fullName:   "r12",
			dwordName:  "r12d",
			wordName:   "r12w",
			byteName:   "r12b",
		},
		R13: {
			integer:    true,
			callerSave: false,
			fullName:   "r13",
			dwordName:  "r13d",
			wordName:   "r13w",
			byteName:   "r13b",
		},
		R14: {
			integer:    true,
			callerSave: false,
			fullName:   "r14",
			dwordName:  "r14d",
			wordName:   "r14w",
			byteName:   "r14b",
		},
		R15: {
			integer:    true,
			callerSave: false,
			fullName:   "r15",
			dwordName:  "r15d",
			wordName:   "r15w",
			byteName:   "r15b",
		},
	}

	return registers
}
