package common

import (
	"crypto/sha1"
	"encoding/hex"
)

func HashsumSHA1(val string) string {
	hashBytes := sha1.Sum([]byte(val))
	return hex.EncodeToString(hashBytes[:])
}
