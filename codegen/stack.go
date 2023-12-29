package codegen

import "sort"

// A Stack represents the use of the stack based relative to RBP.
type Stack struct {
	end    int   // positive offset (16-byte aligned) of start of unused space
	free4  []int // sorted offsets of 4-byte blocks
	free8  []int // sorted offsets of 8-byte blocks
	free16 []int // sorted offsets of 16-byte blocks
	free32 []int // offsets of 16-byte blocks
}

func NewStack() *Stack {
	return &Stack{}
}

// Alloc4 allocates a 4-byte block.
func (s *Stack) Alloc4() int {
	var offset int
	free4 := func(offset int) []int {
		s.Free4(offset)
		return s.free4
	}
	offset, s.free4 = allocN(4, s.free4, s.Alloc8, free4)
	return offset
}

// Free4 deallocates a 4-byte block.
func (s *Stack) Free4(offset int) {
	s.free4 = freeN(s.free4, offset, 4, s.Free8)
}

// Alloc8 allocates an 8-byte block.
func (s *Stack) Alloc8() int {
	var offset int
	free8 := func(offset int) []int {
		s.Free8(offset)
		return s.free8
	}
	offset, s.free8 = allocN(8, s.free8, s.Alloc16, free8)
	return offset
}

// Free8 deallocates a 8-byte block.
func (s *Stack) Free8(offset int) {
	s.free8 = freeN(s.free8, offset, 8, s.Free16)
}

// Alloc16 allocates an 16-byte block.
func (s *Stack) Alloc16() int {
	var offset int
	free16 := func(offset int) []int {
		s.Free16(offset)
		return s.free16
	}
	offset, s.free16 = allocN(16, s.free16, s.alloc32_, free16)
	return offset
}

// Free16 deallocates a 16-byte block.
func (s *Stack) Free16(offset int) {
	s.free16 = freeN(s.free16, offset, 16, s.free32_)
}

func (s *Stack) alloc32_() int {
	var offset int
	if len(s.free32) > 0 {
		offset = s.free32[0]
		s.free32 = s.free32[1:]
	} else {
		offset = s.end
		s.end += 32
	}
	return offset
}

func (s *Stack) free32_(offset int) {
	s.free32 = append(s.free32, offset)
}

func allocN(n int, free []int, alloc2N func() int, freeN func(int) []int) (int, []int) {
	var offset int
	if len(free) > 0 {
		offset = free[0]
		free = free[1:]
	} else {
		offset = alloc2N()
		free = freeN(offset + n)
	}
	return offset, free
}

func freeN(free []int, offset, n int, fNext func(int)) []int {
	free = append(free, offset)
	sort.Ints(free)

	// merge contiguous (2n)-blocks
	n2 := 2 * n
	for i := 1; i < len(free); i++ {
		offset1 := free[i-1] // previous value
		offset2 := free[i]   // current value

		if (offset1&n2) == 0 && (offset1+n) == offset2 {
			// values are contiguous, promote them together
			fNext(offset1)

			// remove them from our slice
			s1 := free[0 : i-1]
			s2 := free[i+1:]
			free = append(s1, s2...)

			// only one merge is ever possible
			break
		}
	}

	return free
}
