//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"

	wire "github.com/buildry-wire/wire-go"
)

func main() {
	client := wire.NewClient("sk_live_...")
	ctx := context.Background()

	pi, err := client.PaymentIntents.Create(ctx, &wire.PaymentIntentCreateParams{Amount: 50000, Currency: "MNT"})
	if err != nil {
		log.Fatal(err)
	}
	confirmed, err := client.PaymentIntents.Confirm(ctx, pi.ID, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(confirmed.ID, confirmed.Status)
}
