package auth

import (
	"fmt"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	fmt.Println(PasswordEncrypt("Password", "Jf7ZgvbH5TfH"))
}
