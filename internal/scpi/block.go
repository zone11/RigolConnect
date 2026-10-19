package scpi

import (
	"fmt"
	"io"
)

// ReadBlock reads one IEEE-488.2 definite-length block, leaving any terminator
// unread. maxBytes limits allocation before any payload is read. Indefinite
// blocks (#0) are unsupported because binary data has no safe line delimiter.
func ReadBlock(r io.Reader, maxBytes int) ([]byte, error) {
	if maxBytes <= 0 {
		return nil, fmt.Errorf("block size limit must be positive")
	}
	var header [2]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, fmt.Errorf("read block header: %w", err)
	}
	if header[0] != '#' || header[1] < '1' || header[1] > '9' {
		return nil, fmt.Errorf("invalid or unsupported block header %q", header[:])
	}
	var digits [9]byte
	width := int(header[1] - '0')
	if _, err := io.ReadFull(r, digits[:width]); err != nil {
		return nil, fmt.Errorf("read block length: %w", err)
	}
	var length int64
	for _, digit := range digits[:width] {
		if digit < '0' || digit > '9' {
			return nil, fmt.Errorf("invalid block length %q", digits[:width])
		}
		length = length*10 + int64(digit-'0')
	}
	if length > int64(maxBytes) {
		return nil, fmt.Errorf("block length %d exceeds limit %d", length, maxBytes)
	}
	payload := make([]byte, int(length))
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, fmt.Errorf("read %d-byte block payload: %w", length, err)
	}
	return payload, nil
}
