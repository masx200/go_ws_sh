package password_hashed

import (
	"strings"
	"testing"
)

func TestHashPasswordWithSaltSHA256(t *testing.T) {
	options := Options{
		algorithm:  "SHA-256",
		saltLength: 32,
	}
	result, err := HashPasswordWithSalt("pass", options)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result.String())
	if result.Algorithm != "SHA-256" {
		t.Errorf("Expected algorithm to be SHA-256, got %s", result.Algorithm)
	}
}

func TestHashPasswordWithSaltSHA384(t *testing.T) {
	options := Options{
		algorithm:  "SHA-384",
		saltLength: 48,
	}
	result, err := HashPasswordWithSalt("pass", options)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result.String())
	if result.Algorithm != "SHA-384" {
		t.Errorf("Expected algorithm to be SHA-384, got %s", result.Algorithm)
	}
}

func TestHashPasswordWithSaltDefault(t *testing.T) {
	result, err := HashPasswordWithSalt("pass")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result.String())
	if result.Algorithm != "SHA-512" {
		t.Errorf("Expected algorithm to be SHA-512, got %s", result.Algorithm)
	}
}

func TestHashPasswordWithSaltSHA512(t *testing.T) {
	result, err := HashPasswordWithSalt("pass")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result.String())
	if result.Algorithm != "SHA-512" {
		t.Errorf("Expected algorithm to be SHA-512, got %s", result.Algorithm)
	}
}

func TestHashPasswordWithSaltShouldUseProvidedSaltHex(t *testing.T) {
	options := Options{
		algorithm:  "SHA-512",
		saltHex:    strings.Repeat("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", 2),
		saltLength: 64,
	}
	result, err := HashPasswordWithSalt("password", options)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result.String())
	if result.Algorithm != "SHA-512" {
		t.Errorf("Expected algorithm to be SHA-512, got %s", result.Algorithm)
	}
	if result.Salt != strings.Repeat("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", 2) {
		t.Errorf("Expected salt to be %s, got %s", strings.Repeat("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", 2), result.Salt)
	}
}
