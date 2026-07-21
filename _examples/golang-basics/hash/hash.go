// Package hash provides a minimal, dependency-free type used to demonstrate
// the `go.type` / `go.type.import` RIDL annotations pointing at a package
// local to this repository, instead of a third-party module.
package hash

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// Hash is a fixed-size, hex-encoded 32-byte hash.
type Hash [32]byte

func (h Hash) String() string {
	return "0x" + hex.EncodeToString(h[:])
}

func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *Hash) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	s = strings.TrimPrefix(s, "0x")
	b, err := hex.DecodeString(s)
	if err != nil {
		return fmt.Errorf("hash: invalid hex string: %w", err)
	}
	if len(b) != len(h) {
		return fmt.Errorf("hash: invalid length %d, expected %d", len(b), len(h))
	}

	copy(h[:], b)
	return nil
}
