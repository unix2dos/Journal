package utils

import (
	"encoding/base64"
	"log"

	"golang.org/x/crypto/scrypt"
)

func StringContains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func IntContains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func BoolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func SliceRemove(slice []int64, index int) []int64 {
	return append(slice[:index], slice[index+1:]...)
}

func SliceRemoveValue(slice []int64, value int64) []int64 {
	for k, v := range slice {
		if value == v {
			return SliceRemove(slice, k)
		}
	}
	return slice
}

func ScryptPassWord(pass string) string {
	salt := []byte{0xc1, 0x08, 0xf2, 0x58, 0xa7, 0x6a, 0xad, 0x7b}
	dk, err := scrypt.Key([]byte(pass), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(dk)
}
