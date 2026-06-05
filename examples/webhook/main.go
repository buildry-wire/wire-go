//go:build ignore

package main

import (
	"fmt"
	"io"
	"net/http"

	wire "github.com/buildry-wire/wire-go"
)

func main() {
	client := wire.NewClient("sk_live_...")
	secret := "whsec_..."

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		event, err := client.Webhooks.Verify(body, r.Header.Get(wire.SignatureHeader), secret)
		if err != nil {
			http.Error(w, "bad signature", http.StatusBadRequest)
			return
		}
		fmt.Println("received:", event.Type)
		w.WriteHeader(http.StatusOK)
	})
	_ = http.ListenAndServe(":4242", nil)
}
