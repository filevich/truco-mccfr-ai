package info

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"hash/adler32"
	"log"
)

func ParseHashFn(ID string) func() hash.Hash {
	switch ID {
	case "adler32":
		return func() hash.Hash { return adler32.New() }
	case "sha160":
		return func() hash.Hash { return sha1.New() }
	case "sha256":
		return func() hash.Hash { return sha256.New() }
	case "sha512":
		return func() hash.Hash { return sha512.New() }
	}
	log.Panicf("hash `%s` not found", ID)
	return nil
}
