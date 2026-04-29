package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

func mainGenHash() {
	password := "19&2H1)a"
	var sha512Hasher = sha512.New()
	sha512Hasher.Write([]byte(password))
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)
	fmt.Println(hashedPasswordHex)
}
