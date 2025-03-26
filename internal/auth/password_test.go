package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	// positive test - correct password
	err = CheckPasswordHash(hashedPassword, password)
	if err != nil {
		t.Errorf("mismatch between password and hash: %v", err)
	}

	// negative test - wrong pasword
	wrongPassword := "wrong-password"
	err = CheckPasswordHash(hashedPassword, wrongPassword)
	if err == nil {
		t.Error("checkPassWord did not return an error for incorrect password")
	}

	hashAgain, _ := HashPassword(password)
	if hashAgain == hashedPassword {
		t.Error("hashing the same pasword twice produced idential hashes")
	}
}
