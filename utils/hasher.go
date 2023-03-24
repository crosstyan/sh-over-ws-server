package utils

import (
	"bytes"

	"github.com/google/uuid"
)

type UuidHasher struct{}

// A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC
// 4122, which is way bigger than uint32.
func (h *UuidHasher) Hash(key uuid.UUID) uint32 {
	return key.ID()
}

func (h *UuidHasher) Equal(a, b uuid.UUID) bool {
	return bytes.Equal(a[:], b[:])
}
