// pkg/cbor/codec.go
package cbor

import (
	"github.com/fxamacker/cbor/v2"
)

// encMode holds the configured encoding options.
// We use Canonical encoding to ensure deterministic output (sorted map keys).
// This is critical for generating consistent checksums for configuration storage.
var encMode cbor.EncMode

func init() {
	var err error
	// Initialize the encoder options.
	// We panic here if initialization fails because the framework depends on this
	// for the ConfigManager. A failure here is unrecoverable.
	encMode, err = cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		panic("cbor: codec init failed: " + err.Error())
	}
}

// Marshal serializes a Go value into CBOR format.
//
// Performance Note: This function uses runtime reflection.
// In TinyGo, this will increase binary size. Minimize usage to
// configuration saving/loading and avoid using it in tight loops.
func Marshal(v interface{}) ([]byte, error) {
	return encMode.Marshal(v)
}

// Unmarshal deserializes CBOR data into a Go value.
// 'v' must be a non-nil pointer.
func Unmarshal(data []byte, v interface{}) error {
	return cbor.Unmarshal(data, v)
}