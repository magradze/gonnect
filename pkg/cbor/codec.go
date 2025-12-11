// pkg/cbor/codec.go
package cbor

import (
	"github.com/fxamacker/cbor/v2"
)

// EncMode defines the canonical encoding options.
// We use Canonical encoding to ensure deterministic output (e.g., map keys are sorted).
// This is crucial for comparing configurations or generating checksums in embedded systems.
var EncMode, _ = cbor.CanonicalEncOptions().EncMode()

// Marshal serializes a Go value into CBOR format.
// It acts as a lightweight wrapper around the underlying CBOR library,
// abstracting the specific implementation details from the rest of the framework.
func Marshal(v interface{}) ([]byte, error) {
	return EncMode.Marshal(v)
}

// Unmarshal deserializes CBOR data into a Go value.
// The target value 'v' must be a non-nil pointer.
func Unmarshal(data []byte, v interface{}) error {
	return cbor.Unmarshal(data, v)
}