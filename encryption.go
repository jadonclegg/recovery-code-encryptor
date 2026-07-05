package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

func EncryptBase64(req *encryptionRequest) ([]byte, error) {
	key := getCipherKey(req.Password, req.Name)
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher: %w", err)
	}

	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(req.Data)+aead.Overhead())
	rand.Read(nonce)

	cypherText := aead.Seal(nonce, nonce, []byte(req.Data), []byte(req.Name))

	out := make([]byte, 0, len(req.Data)*2)
	buf := bytes.NewBuffer(out)
	enc := base64.NewEncoder(base64.StdEncoding, buf)
	_, err = enc.Write(cypherText)
	if err != nil {
		return nil, fmt.Errorf("error converting to base64: %w", err)
	}

	return buf.Bytes(), nil
}

func DecryptBase64(req *encryptionRequest) ([]byte, error) {
	key := getCipherKey(req.Password, req.Name)
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher: %w", err)
	}

	encryptedMsg, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64 string: %w", err)
	}

	if len(encryptedMsg) < aead.NonceSize() {
		err = errors.New("ciphertext too short")
		return nil, err
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := encryptedMsg[:aead.NonceSize()], encryptedMsg[aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintext, err := aead.Open(nil, nonce, ciphertext, []byte(req.Name))
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %w", err)
	}

	return plaintext, nil
}
