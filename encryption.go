package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

const SALT_LEN = 16

func EncryptBase64(req *encryptionRequest) ([]byte, error) {
	key := getCipherKeyString(req.Password, req.Name)
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher: %w", err)
	}

	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(req.Data)+aead.Overhead())
	rand.Read(nonce)

	cypherText := aead.Seal(nonce, nonce, []byte(req.Data), []byte(req.Name))

	enc := base64.StdEncoding.EncodeToString(cypherText)

	return []byte(enc), nil
}

func EncryptBase64V2(req *encryptionRequest) ([]byte, error) {
	salt := make([]byte, SALT_LEN)
	rand.Read(salt)

	key := getCipherKey([]byte(req.Password), salt)

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher: %w", err)
	}

	out := make([]byte, SALT_LEN+aead.NonceSize(), aead.NonceSize()+len(req.Data)+aead.Overhead())
	copy(out, salt)

	nonce := out[SALT_LEN:]
	rand.Read(nonce)

	cypherText := aead.Seal(out, nonce, []byte(req.Data), []byte(req.Name))

	enc := base64.StdEncoding.EncodeToString(cypherText)

	return []byte(enc), nil
}

func DecryptBase64(req *encryptionRequest) ([]byte, error) {
	key := getCipherKeyString(req.Password, req.Name)
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

func DecryptBase64V2(req *encryptionRequest) ([]byte, error) {
	encryptedMsg, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64 string: %w", err)
	}

	if len(encryptedMsg) < SALT_LEN {
		return nil, fmt.Errorf("ciphertext too short")
	}

	salt := encryptedMsg[:SALT_LEN]

	key := getCipherKey([]byte(req.Password), salt)
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher: %w", err)
	}

	if len(encryptedMsg) < aead.NonceSize()+SALT_LEN {
		err = errors.New("ciphertext too short")
		return nil, err
	}

	encryptedMsg = encryptedMsg[SALT_LEN:]
	// Split nonce and ciphertext.
	nonce, ciphertext := encryptedMsg[:aead.NonceSize()], encryptedMsg[aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintext, err := aead.Open(nil, nonce, ciphertext, []byte(req.Name))
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %w", err)
	}

	return plaintext, nil
}

func getCipherKeyString(password string, salt string) []byte {
	return getCipherKey([]byte(password), []byte(salt))
}

func getCipherKey(password []byte, salt []byte) []byte {
	return argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
}
