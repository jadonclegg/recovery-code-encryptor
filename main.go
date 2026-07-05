package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	_ "embed"

	"golang.org/x/crypto/argon2"
)

//go:embed web2/dist/web2/browser
var web embed.FS

func main() {
	r := http.NewServeMux()

	r.HandleFunc("/encrypt", handleEncrypt)

	r.HandleFunc("/decrypt", handleDecrypt)

	if len(os.Args) > 1 && os.Args[1] == "--fs" {
		fmt.Println("using fs mode")
		r.Handle("/", http.FileServer(http.Dir("web2/dist/web2/browser")))
	} else {
		fSys, err := fs.Sub(web, "web2/dist/web2/browser")
		if err != nil {
			panic(err)
		}

		r.Handle("/", http.FileServer(http.FS(fSys)))
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

	ciphertext, err := EncryptBase64(req)
	if err != nil {
		handleEncryptError(w, response, err)
		return
	}

	response.Result = string(ciphertext)

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

	plaintext, err := DecryptBase64(req)
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
