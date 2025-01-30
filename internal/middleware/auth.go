package auth

import (
	"crypto/sha256"
)

var PassCodeHash [32]byte = sha256.Sum256([]byte("12345"))
