package codegen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStack1(t *testing.T) {
	r := require.New(t)
	s := NewStack()

	// allocate a value
	offset := s.Alloc4()
	r.Equal(0, offset)
	r.Equal(1, len(s.free4))
	r.Equal(1, len(s.free8))
	r.Equal(1, len(s.free16))
	r.Equal(0, len(s.free32))

	// free it again
	s.Free4(offset)
	r.Equal(0, len(s.free4))
	r.Equal(0, len(s.free8))
	r.Equal(0, len(s.free16))
	r.Equal(1, len(s.free32))
}

func TestStack2(t *testing.T) {
	r := require.New(t)
	s := NewStack()

	// allocate two value
	offset1 := s.Alloc4()
	offset2 := s.Alloc4()
	r.Equal(0, offset1)
	r.Equal(4, offset2)
	r.Equal(0, len(s.free4))
	r.Equal(1, len(s.free8))
	r.Equal(1, len(s.free16))
	r.Equal(0, len(s.free32))

	// free them again
	s.Free4(offset1)
	s.Free4(offset2)
	r.Equal(0, len(s.free4))
	r.Equal(0, len(s.free8))
	r.Equal(0, len(s.free16))
	r.Equal(1, len(s.free32))
}

func TestStack3(t *testing.T) {
	r := require.New(t)
	s := NewStack()

	// allocate two value
	offset1 := s.Alloc4()
	offset2 := s.Alloc4()
	r.Equal(0, offset1)
	r.Equal(4, offset2)
	r.Equal(0, len(s.free4))
	r.Equal(1, len(s.free8))
	r.Equal(1, len(s.free16))
	r.Equal(0, len(s.free32))

	// free them again (opposite order)
	s.Free4(offset2)
	s.Free4(offset1)
	r.Equal(0, len(s.free4))
	r.Equal(0, len(s.free8))
	r.Equal(0, len(s.free16))
	r.Equal(1, len(s.free32))
}
