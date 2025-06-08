// # tools/gen_hash.go
package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	h, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	fmt.Println(string(h))
}
