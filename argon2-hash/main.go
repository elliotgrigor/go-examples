package main

import (
	"crypto/rand"
	"fmt"
)

func main() {
	password := "foobarbaz69"
	salt := rand.Text() + rand.Text()

	fs := argon2idFormatString(password, salt)

	// password = "badpassword"
	ok, err := argon2idCompare(password, fs, salt)
	if err != nil {
		panic(err)
	}

	if ok {
		fmt.Println("Success")
	} else {
		fmt.Println("Fail")
	}
}
