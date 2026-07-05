package main

import (
	"slices"
	"testing"
)

func TestEncryptDecryptWorks(t *testing.T) {
	encReq := &encryptionRequest{
		Name:     "name",
		Data:     "asdf :D",
		Password: "asdfasdf",
	}

	result, err := EncryptBase64(encReq)
	if err != nil {
		t.Errorf("error encrypting base64: %s", err)
		return
	}

	if slices.Equal(result, []byte(encReq.Data)) {
		t.Errorf("result should not be the plaintext!!!")
		return
	}

	decReq := &encryptionRequest{
		Name:     "name",
		Data:     string(result),
		Password: "asdfasdf",
	}

	plain, err := DecryptBase64(decReq)
	if err != nil {
		t.Errorf("error while decrypting: %s", err)
		return
	}

	if !slices.Equal(plain, []byte(encReq.Data)) {
		t.Error("expected result to equal the original data, but was not the same")
		return
	}

	decReq.Name = "a"
	plain, err = DecryptBase64(decReq)
	if err == nil {
		t.Errorf("decryption should have failed but didn't...")
		return
	}

	decReq.Name = "asdf"
	decReq.Password = "asdf"
	plain, err = DecryptBase64(decReq)
	if err == nil {
		t.Errorf("decryption should have failed but didn't...")
		return
	}
}

func TestEncryptDecryptV2Works(t *testing.T) {
	encReq := &encryptionRequest{
		Name:     "name",
		Data:     "asdf :D",
		Password: "asdfasdf",
	}

	result, err := EncryptBase64V2(encReq)
	if err != nil {
		t.Errorf("error encrypting base64: %s", err)
		return
	}

	if slices.Equal(result, []byte(encReq.Data)) {
		t.Errorf("result should not be the plaintext!!!")
		return
	}

	decReq := &encryptionRequest{
		Name:     "name",
		Data:     string(result),
		Password: "asdfasdf",
	}

	plain, err := DecryptBase64V2(decReq)
	if err != nil {
		t.Errorf("error while decrypting: %s", err)
		return
	}

	if !slices.Equal(plain, []byte(encReq.Data)) {
		t.Error("expected result to equal the original data, but was not the same")
		return
	}

	decReq.Name = "a"
	plain, err = DecryptBase64V2(decReq)
	if err == nil {
		t.Errorf("decryption should have failed but didn't...")
		return
	}

	decReq.Name = "asdf"
	decReq.Password = "asdf"
	plain, err = DecryptBase64V2(decReq)
	if err == nil {
		t.Errorf("decryption should have failed but didn't...")
		return
	}
}
