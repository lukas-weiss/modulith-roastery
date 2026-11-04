// Package uuid makes random (version 4) UUIDs, like java.util.UUID.randomUUID().
// The standard library has none, and the example should need no dependencies.
package uuid

import (
	"crypto/rand"
	"fmt"
)

// New returns a random UUID such as "1b4e28ba-2fa1-41d2-883f-0016d3cca427".
func New() string {
	var b [16]byte
	_, _ = rand.Read(b[:])  // never fails, see crypto/rand
	b[6] = b[6]&0x0f | 0x40 // version 4
	b[8] = b[8]&0x3f | 0x80 // variant RFC 4122
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
