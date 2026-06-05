# wire-go

Official Go SDK for the [Wire](https://wire.mn) payment API.

## Install
```bash
go get github.com/buildry-wire/wire-go
```

## Quickstart
```go
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

	pi, err := client.PaymentIntents.Create(ctx, &wire.PaymentIntentCreateParams{
		Amount:   50000, // 500.00 MNT, minor units
		Currency: "MNT",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pi.ID, pi.Status)
}
```

## Auto-pagination
```go
it := client.Charges.List(ctx, &wire.ListParams{Limit: 50})
for it.Next() {
	fmt.Println(it.Current().ID)
}
if err := it.Err(); err != nil {
	log.Fatal(err)
}
```

## Webhook verification
```go
func handler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	event, err := client.Webhooks.Verify(body, r.Header.Get(wire.SignatureHeader), endpointSecret)
	if err != nil {
		http.Error(w, "bad signature", http.StatusBadRequest)
		return
	}
	fmt.Println("event:", event.Type)
}
```

## Errors
```go
var werr *wire.Error
if errors.As(err, &werr) {
	fmt.Println(werr.Code, werr.RequestID, werr.StatusCode)
}
```

## License
MIT
