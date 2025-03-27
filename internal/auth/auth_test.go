package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
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
func TestValidateJWT(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, "secret", time.Hour)
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	// valid token so passes
	id, err := ValidateJWT(token, "secret")
	if err != nil {
		t.Errorf("Error validating token: %v", err)
	}
	if userID != id {
		t.Errorf("Actual id: %s\nExpected id: %s", id, userID)
	}

	// wrong token so does not pass
	_, err = ValidateJWT("wrongToken", "secret")
	if err == nil {
		t.Error("Expected validation to fail for incorrect token")
	}

	// wronger secret so does not pass
	_, err = ValidateJWT(token, "wrongSecret")
	if err == nil {
		t.Error("Expected validation to fail for incorrect secret")
	}

	// expired token so does not pass
	expiredToken, err := MakeJWT(userID, "secret", -1*time.Hour)
	if err != nil {
		t.Errorf("error creating expired JWT: %v", err)
	}
	_, err = ValidateJWT(expiredToken, "secret")
	if err == nil {
		t.Error("Expected validation to fail for expired token")
	}
}

func TestGetBearerToken(t *testing.T) {
	userID := uuid.New()
	tokenString, err := MakeJWT(userID, "secret", time.Hour)
	if err != nil {
		t.Error("error creating JWT:")
	}

	emptyHeader := http.Header{}
	_, err = GetBearerToken(emptyHeader)
	if err == nil {
		t.Errorf("Authorization not set in http header: %v", err)
	}

	invalidHeader := http.Header{}
	invalidHeader.Set("Authorization", tokenString)
	_, err = GetBearerToken(invalidHeader)
	if err == nil {
		t.Error("expected error parsing http header without Bearing prefix")
	}

	validHeader := http.Header{}
	validHeader.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	token, err := GetBearerToken(validHeader)
	if err != nil {
		t.Errorf("error parsing http header: %v", err)
	}
	if token != tokenString {
		t.Errorf("token mismatch: got %s, expected %s", token, tokenString)
	}
}
