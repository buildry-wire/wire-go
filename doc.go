// Package wire is the official Go SDK for the Wire payment API (https://wire.mn).
//
// Create a client with an API key and call resource services:
//
//	client := wire.NewClient("sk_live_...")
//	pi, err := client.PaymentIntents.Create(ctx, &wire.PaymentIntentCreateParams{Amount: 50000})
package wire
