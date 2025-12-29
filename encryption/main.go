package main

import "fmt"

func main() {
	key := []byte("du9CHhLDwmk7ZdwI0AV8DpM3UqbgQGvk") // must be 32 bytes
	message := []byte("urmom")

	enc, err := Encrypt(key, message)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(enc))

	orig, err := Decrypt(key, enc)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(orig))
}
