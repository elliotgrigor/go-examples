package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	memoryCost  = 256 * 1024 // 256MB
	timeFactor  = 10
	parallelism = 6
)

func argon2idFormatString(text, salt string) string {
	hash := argon2.IDKey([]byte(text), []byte(salt), timeFactor, memoryCost, parallelism, uint32(len(salt)))
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	sb := strings.Builder{}
	fmt.Fprintf(&sb, "$argon2id$v=19$m=%d", memoryCost)
	fmt.Fprintf(&sb, ",t=%d", timeFactor)
	fmt.Fprintf(&sb, ",p=%d", parallelism)
	fmt.Fprintf(&sb, "$%s", hashBase64)

	return sb.String()
}

func argon2idExtractBase64Hash(fmtStr string) ([]byte, error) {
	components := strings.Split(fmtStr, "$")
	if len(components) != 5 {
		return []byte{}, errors.New("argon2idExtractBase64Hash: incorrect number of components")
	}
	hashBase64 := components[4]
	hash, err := base64.StdEncoding.DecodeString(hashBase64)
	if err != nil {
		return []byte{}, fmt.Errorf("argon2idExtractBase64Hash: %s", err.Error())
	}
	return hash, nil
}

func argon2idCompare(compare, fmtStr, salt string) (bool, error) {
	tryHash := argon2.IDKey([]byte(compare), []byte(salt), timeFactor, memoryCost, parallelism, uint32(len(salt)))
	hash, err := argon2idExtractBase64Hash(fmtStr)
	if err != nil {
		return false, err
	}
	if string(tryHash) == string(hash) {
		return true, nil
	}
	return false, nil
}
