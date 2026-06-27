package main

import (
	"crypto/rand"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	_ "embed"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

//go:embed web2/dist/web2/browser
var web embed.FS

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/encrypt", handleEncrypt)

	r.HandleFunc("/decrypt", handleDecrypt)

	if len(os.Args) > 1 && os.Args[1] == "--fs" {
		fmt.Println("using fs mode")
		r.PathPrefix("/").Handler(http.FileServer(http.Dir("dist/web2/browser")))
	} else {
		fSys, err := fs.Sub(web, "web2/dist/web2/browser")
		if err != nil {
			panic(err)
		}

		r.PathPrefix("/").Handler(http.FileServer(http.FS(fSys)))
	}

	http.Handle("/", r)

	fmt.Println("http server listening on :8080")

	err := http.ListenAndServe(":8080", nil)
	panic(err)
}

type encryptionResponse struct {
	Result  string `json:"result"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type encryptionRequest struct {
	Password string `json:"password"`
	Name     string `json:"name"`
	Data     string `json:"data"`
}

func handleEncrypt(w http.ResponseWriter, r *http.Request) {
	response := &encryptionResponse{
		Result:  "result",
		Success: true,
		Message: "",
	}

	dec := json.NewDecoder(r.Body)
	req := &encryptionRequest{}
	err := dec.Decode(req)
	if err != nil {
		response.Success = false
		response.Message = "failed to decode json"
		writeEncrpytionResponse(w, response)
		return
	}

	key := getCipherKey(req.Password, req.Name)
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		handleEncryptError(w, response, err)
		return
	}

	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(req.Data)+aead.Overhead())
	_, err = rand.Read(nonce)
	if err != nil {
		handleEncryptError(w, response, err)
		return
	}

	cypherText := aead.Seal(nonce, nonce, []byte(req.Data), []byte(req.Name))

	response.Result = base64.StdEncoding.EncodeToString(cypherText)
	writeEncrpytionResponse(w, response)
}

func handleEncryptError(w http.ResponseWriter, response *encryptionResponse, err error) {
	response.Success = false
	response.Message = err.Error()
	writeEncrpytionResponse(w, response)
}

func writeEncrpytionResponse(w http.ResponseWriter, response *encryptionResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to marshal response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

func handleDecrypt(w http.ResponseWriter, r *http.Request) {
	response := &encryptionResponse{
		Result:  "result",
		Success: true,
		Message: "",
	}

	dec := json.NewDecoder(r.Body)
	req := &encryptionRequest{}
	err := dec.Decode(req)
	if err != nil {
		response.Success = false
		response.Message = "failed to decode json"
		writeEncrpytionResponse(w, response)
		return
	}

	key := getCipherKey(req.Password, req.Name)
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		handleEncryptError(w, response, err)
		return
	}

	encryptedMsg, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		response.Success = false
		response.Message = "failed to decode base64"
		writeEncrpytionResponse(w, response)
		return
	}

	if len(encryptedMsg) < aead.NonceSize() {
		err = errors.New("ciphertext too short")
		handleEncryptError(w, response, err)
		return
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := encryptedMsg[:aead.NonceSize()], encryptedMsg[aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintext, err := aead.Open(nil, nonce, ciphertext, []byte(req.Name))
	if err != nil {
		handleEncryptError(w, response, err)
		return
	}

	response.Result = string(plaintext)
	writeEncrpytionResponse(w, response)
}

func getCipherKey(password string, salt string) []byte {
	return argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
}
