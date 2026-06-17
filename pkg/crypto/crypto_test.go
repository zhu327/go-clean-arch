package crypto

import (
	"testing"
)

func TestBcryptHasher_Hash(t *testing.T) {
	h := NewBcryptHasher()
	hash, err := h.Hash("secret")
	if err != nil {
		t.Fatalf("Hash returned error: %v", err)
	}
	if hash == "" {
		t.Fatal("Hash returned empty string")
	}
	if hash == "secret" {
		t.Fatal("Hash returned plaintext password")
	}
}

func TestBcryptHasher_Verify(t *testing.T) {
	h := NewBcryptHasher()
	hash, err := h.Hash("secret")
	if err != nil {
		t.Fatalf("Hash returned error: %v", err)
	}
	if err := h.Verify("secret", hash); err != nil {
		t.Fatalf("Verify correct password returned error: %v", err)
	}
	if err := h.Verify("wrong", hash); err == nil {
		t.Fatal("Verify wrong password expected error, got nil")
	}
}
