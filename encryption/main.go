package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// Setup shared between encryption and decryption
func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm, nil
}

func encrypt(key, plainText []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plainText, nil), nil
}

func decrypt(key, cipherText []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce := cipherText[:nonceSize]
	data := cipherText[nonceSize:]
	return gcm.Open(nil, nonce, data, nil)
}

func main() {
	key := []byte("du9CHhLDwmk7ZdwI0AV8DpM3UqbgQGvk") // must be 32 bytes
	message := []byte("urmom")

	enc, err := encrypt(key, message)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(enc))

	orig, err := decrypt(key, enc)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(orig))
}
