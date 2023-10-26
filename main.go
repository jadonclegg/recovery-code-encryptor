package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	_ "embed"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/argon2"
)

//go:embed web/dist/web
var web embed.FS

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/encrypt", handleEncrypt)

	r.HandleFunc("/decrypt", handleDecrypt)

	if len(os.Args) > 1 && os.Args[1] == "--fs" {
		fmt.Println("using fs mode")
		r.PathPrefix("/").Handler(http.FileServer(http.Dir("dist/web")))
	} else {
		fSys, err := fs.Sub(web, "web/dist/web")
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

func handleEncrypt(w http.ResponseWriter, r *http.Request) {

}

func handleDecrypt(w http.ResponseWriter, r *http.Request) {

}

func getCipherKey(password string, salt string) []byte {
	return argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
}
