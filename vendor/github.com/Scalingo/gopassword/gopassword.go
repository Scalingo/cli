package gopassword

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

func Generate(n ...int) string {
	length := 20
	if len(n) > 0 {
		length = n[0]
	}

	// With base64 encoding we need 3 random bytes to
	// have 4 random characters
	randSize := 3 * (length/4 + 1)
	randBytes := make([]byte, randSize)
	rand.Read(randBytes)

	// Encode them in base64
	randString := base64.StdEncoding.EncodeToString(randBytes)

	password := randString[:length]
	password = strings.Replace(password, "+", "_", -1)
	password = strings.Replace(password, "/", "-", -1)

	if password[0] == '-' {
		password = fmt.Sprintf("_%s", password[1:])
	}
	if password[length-1] == '-' {
		password = fmt.Sprintf("%s_", password[:length-1])
	}

	return password
}
