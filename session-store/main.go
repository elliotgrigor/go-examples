package main

import "fmt"

func main() {
	sessions := SessionStore{
		Store: map[string][]byte{},
	}

	sessions.Set("jsmith", SessionData{Subject: "e284359", Iat: "1709055600", Exp: "1709065000"})
	sessions.Set("jdoe", SessionData{Subject: "e284982", Iat: "1709055600", Exp: "1709065000"})
	sessions.Set("jbloggs", SessionData{Subject: "e284123", Iat: "1709055600", Exp: "1709065000"})

	sessions.Delete("jdoe")

	for key := range sessions.Store {
		sess, err := sessions.Get(key)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%-12s %#v\n", key, sess)
	}
}
