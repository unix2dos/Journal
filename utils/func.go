package utils

import (
	"encoding/base64"
	"log"
	"strings"

	"golang.org/x/crypto/scrypt"
)

func StringContains(str string, arr []string) bool {
	for _, v := range arr {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}

func ScryptPassWord(pass string) string {
	salt := []byte{0xc1, 0x08, 0xf2, 0x58, 0xa7, 0x6a, 0xad, 0x7b}
	dk, err := scrypt.Key([]byte(pass), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(dk)
}
